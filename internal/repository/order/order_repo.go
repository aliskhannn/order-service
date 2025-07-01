package order

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	customerr "github.com/aliskhannn/order-service/internal/errors"
	"github.com/aliskhannn/order-service/internal/model"
)

type Repository struct {
	logger *zap.Logger
	db     *pgxpool.Pool
}

func New(l *zap.Logger, db *pgxpool.Pool) *Repository {
	return &Repository{
		logger: l,
		db:     db,
	}
}

func (r *Repository) SaveOrder(ctx context.Context, order *model.Order) (string, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return "", fmt.Errorf("begin tx: %w", customerr.ErrTxBegin)
	}
	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				r.logger.Error("failed to rollback transaction", zap.Error(rollbackErr))
			} else {
				r.logger.Warn("transaction rollback due to error", zap.Error(err))
			}

			return
		}

		if commitErr := tx.Commit(ctx); commitErr != nil {
			r.logger.Error("failed to commit transaction", zap.Error(commitErr))
			err = fmt.Errorf("commit tx: %w", customerr.ErrTxCommit)
		}
	}()

	// Insert into orders
	orderQuery := `
	INSERT INTO orders (
		order_uid, track_number, entry, locale, internal_signature, customer_id,
		delivery_service, shardkey, sm_id, date_created, oof_shard
	) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`

	_, err = tx.Exec(ctx, orderQuery,
		order.OrderID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard,
	)
	if err != nil {
		r.logger.Error("failed to insert order", zap.Error(err), zap.String("order_id", order.OrderID))
		return "", fmt.Errorf("insert order: %w", customerr.ErrInsertOrder)
	}

	// Insert into delivery
	d := order.Delivery
	deliveryQuery := `
	INSERT INTO delivery (
	    order_uid, name, phone, zip, city, address, region, email
	) VALUES($1, $2, $3, $4, $5, $6, $7, $8);
	`

	_, err = tx.Exec(ctx, deliveryQuery,
		order.OrderID, d.Name, d.Phone, d.Zip, d.City, d.Address, d.Region, d.Email)
	if err != nil {
		r.logger.Error("failed to insert delivery", zap.Error(err), zap.String("order_id", order.OrderID))
		return "", fmt.Errorf("insert delivery: %w", customerr.ErrInsertDelivery)
	}

	// Insert into payment
	p := order.Payment
	paymentQuery := `
		INSERT INTO payment (
			transaction, order_uid, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`
	_, err = tx.Exec(ctx, paymentQuery,
		p.Transaction, order.OrderID, p.RequestID, p.Currency, p.Provider,
		p.Amount, p.PaymentDT, p.Bank, p.DeliveryCost, p.GoodsTotal, p.CustomFee)
	if err != nil {
		r.logger.Error("failed to insert payment", zap.Error(err), zap.String("order_id", order.OrderID))
		return "", fmt.Errorf("insert payment: %w", customerr.ErrInsertPayment)
	}

	// Insert each item into items
	itemQuery := `
		INSERT INTO items (
			order_id, chrt_id, track_number, price, rid,
			name, sale, size, total_price, nm_id, brand, status
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
	`
	for _, item := range order.Items {
		_, err = tx.Exec(ctx, itemQuery,
			order.OrderID, item.ChrtID, item.TrackNumber, item.Price, item.RID,
			item.Name, item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			r.logger.Error("failed to insert item", zap.Error(err), zap.String("order_id", order.OrderID))
			return "", fmt.Errorf("insert item: %w", customerr.ErrInsertItem)
		}
	}

	return order.OrderID, nil
}

func (r *Repository) GetOrderById(ctx context.Context, orderID string) (*model.Order, error) {
	var order model.Order

	// Getting order
	orderQuery := `
	SELECT order_uid, track_number, entry, locale, internal_signature, customer_id,
		delivery_service, shardkey, sm_id, date_created, oof_shard
	FROM orders
	WHERE order_uid = $1;
	`

	err := r.db.QueryRow(ctx, orderQuery, orderID).Scan(
		&order.OrderID, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerId, &order.DeliveryService,
		&order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn("order not found", zap.String("order_id", orderID))
			return nil, fmt.Errorf("order with id %s not found: %w", orderID, customerr.ErrOrderNotFound)
		}
		r.logger.Error("failed to get order", zap.Error(err), zap.String("order_id", orderID))
		return nil, fmt.Errorf("error getting order %s: %v", orderID, err)
	}

	// Getting delivery
	deliveryQuery := `
	SELECT name, phone, zip, city, address, region, email
	FROM delivery
    WHERE order_uid = $1;
	`

	var delivery model.Delivery
	err = r.db.QueryRow(ctx, deliveryQuery, orderID).Scan(
		&delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City,
		&delivery.Address, &delivery.Region, &delivery.Email,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn("delivery not found", zap.String("order_id", orderID))
			return nil, fmt.Errorf("delivery for order %s not found: %w", orderID, customerr.ErrDeliveryNotFound)
		}
		r.logger.Error("failed to get delivery", zap.Error(err), zap.String("order_id", orderID))
		return nil, fmt.Errorf("error getting delivery for order %s: %v", orderID, err)
	}
	order.Delivery = delivery

	// Getting payment
	paymentQuery := `
		SELECT transaction, request_id, currency, provider,
			amount, payment_dt, bank, delivery_cost, goods_total, custom_fee
		FROM payment
		WHERE order_uid = $1;
	`

	var payment model.Payment
	err = r.db.QueryRow(ctx, paymentQuery, orderID).Scan(
		&payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider,
		&payment.Amount, &payment.PaymentDT, &payment.Bank, &payment.DeliveryCost,
		&payment.GoodsTotal, &payment.CustomFee,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Warn("payment not found", zap.String("order_id", orderID))
			return nil, fmt.Errorf("payment for order %s not found: %w", orderID, customerr.ErrPaymentNotFound)
		}
		r.logger.Error("failed to get payment", zap.Error(err), zap.String("order_id", orderID))
		return nil, fmt.Errorf("error getting payment for order %s: %v", orderID, err)
	}
	order.Payment = payment

	// Getting items
	itemsQuery := `
		SELECT chrt_id, track_number, price, rid, name, sale,
		       size, total_price, nm_id, brand, status
		FROM items
		WHERE order_id = $1;
	`
	rows, err := r.db.Query(ctx, itemsQuery, orderID)
	if err != nil {
		r.logger.Error("failed to get items", zap.Error(err), zap.String("order_id", orderID))
		return nil, fmt.Errorf("error getting items for order %s: %v", orderID, err)
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err = rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID,
			&item.Name, &item.Sale, &item.Size, &item.TotalPrice,
			&item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			r.logger.Error("failed to scan item", zap.Error(err), zap.String("order_id", orderID))
			return nil, fmt.Errorf("error scanning item for order %s: %v", orderID, customerr.ErrItemScanFailed)
		}
		items = append(items, item)
	}
	order.Items = items

	return &order, nil
}
