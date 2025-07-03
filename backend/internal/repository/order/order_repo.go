package order

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"

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

func (r *Repository) SaveOrder(ctx context.Context, order *model.Order) (uuid.UUID, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		r.logger.Error("failed to begin transaction", zap.Error(err))
		return uuid.Nil, fmt.Errorf("begin tx: %w", customerr.ErrTxBegin)
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
		track_number, entry, locale, internal_signature, customer_id,
		delivery_service, shardkey, sm_id, oof_shard
	) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)
	RETURNING order_uid;
	`

	err = tx.QueryRow(ctx, orderQuery, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.OofShard,
	).Scan(&order.OrderID)
	if err != nil {
		r.logger.Error("failed to insert order", zap.Error(err), zap.String("order_id", order.OrderID.String()))
		return uuid.Nil, fmt.Errorf("insert order: %w", customerr.ErrInsertOrder)
	}
	r.logger.Info("order inserted", zap.String("order_id", order.OrderID.String()))

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
		r.logger.Error("failed to insert delivery", zap.Error(err), zap.String("order_id", order.OrderID.String()))
		return uuid.Nil, fmt.Errorf("insert delivery: %w", customerr.ErrInsertDelivery)
	}
	r.logger.Info("delivery inserted", zap.String("order_id", order.OrderID.String()))

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
		r.logger.Error("failed to insert payment", zap.Error(err), zap.String("order_id", order.OrderID.String()))
		return uuid.Nil, fmt.Errorf("insert payment: %w", customerr.ErrInsertPayment)
	}
	r.logger.Info("payment inserted", zap.String("order_id", order.OrderID.String()))

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
			r.logger.Error("failed to insert item", zap.Error(err), zap.String("order_id", order.OrderID.String()))
			return uuid.Nil, fmt.Errorf("insert item: %w", customerr.ErrInsertItem)
		}
	}
	r.logger.Info("items inserted", zap.String("order_id", order.OrderID.String()))

	return order.OrderID, nil
}

func (r *Repository) GetOrderById(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	query := `
	SELECT
		o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id,
		o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
	
		d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
	
		p.transaction, p.request_id, p.currency, p.provider,
		p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
	FROM orders o
	JOIN delivery d ON o.order_uid = d.order_uid
	JOIN payment p ON o.order_uid = p.order_uid
	WHERE o.order_uid = $1;
	`

	row := r.db.QueryRow(ctx, query, orderID)

	var o model.Order
	var d model.Delivery
	var p model.Payment

	err := row.Scan(
		// Order
		&o.OrderID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature, &o.CustomerId,
		&o.DeliveryService, &o.Shardkey, &o.SmId, &o.DateCreated, &o.OofShard,

		// Delivery
		&d.Name, &d.Phone, &d.Zip, &d.City, &d.Address, &d.Region, &d.Email,

		// Payment
		&p.Transaction, &p.RequestID, &p.Currency, &p.Provider,
		&p.Amount, &p.PaymentDT, &p.Bank, &p.DeliveryCost, &p.GoodsTotal, &p.CustomFee,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("get order by id: %w", customerr.ErrOrderNotFound)
		}

		return nil, fmt.Errorf("scan row: %w", customerr.ErrScanRow)
	}

	o.Delivery = d
	o.Payment = p

	return &o, err
}

func (r *Repository) GetLastOrders(ctx context.Context, limit int) ([]model.Order, error) {
	query := `
	SELECT
		o.order_uid, o.track_number, o.entry, o.locale, o.internal_signature, o.customer_id,
		o.delivery_service, o.shardkey, o.sm_id, o.date_created, o.oof_shard,
	
		d.name, d.phone, d.zip, d.city, d.address, d.region, d.email,
	
		p.transaction, p.request_id, p.currency, p.provider,
		p.amount, p.payment_dt, p.bank, p.delivery_cost, p.goods_total, p.custom_fee
	FROM orders o
	JOIN delivery d ON o.order_uid = d.order_uid
	JOIN payment p ON o.order_uid = p.order_uid
	ORDER BY o.date_created DESC
	LIMIT $1
	`

	rows, err := r.db.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("get last orders: %w", customerr.ErrGetLastOrders)
	}
	defer rows.Close()

	var orders []model.Order

	for rows.Next() {
		var o model.Order
		var d model.Delivery
		var p model.Payment

		err = rows.Scan(
			// Order
			&o.OrderID, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature, &o.CustomerId,
			&o.DeliveryService, &o.Shardkey, &o.SmId, &o.DateCreated, &o.OofShard,

			// Delivery
			&d.Name, &d.Phone, &d.Zip, &d.City, &d.Address, &d.Region, &d.Email,

			// Payment
			&p.Transaction, &p.RequestID, &p.Currency, &p.Provider,
			&p.Amount, &p.PaymentDT, &p.Bank, &p.DeliveryCost, &p.GoodsTotal, &p.CustomFee,
		)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("get order by id: %w", customerr.ErrOrderNotFound)
			}

			return nil, fmt.Errorf("scan row: %w", customerr.ErrScanRow)
		}

		o.Delivery = d
		o.Payment = p

		items, err := r.GetItemsByOrderID(ctx, o.OrderID)
		if err != nil {
			return nil, fmt.Errorf("get items by order id: %w", customerr.ErrGetItemsByOrderId)
		}
		o.Items = items

		orders = append(orders, o)
	}

	return orders, nil
}

func (r *Repository) GetItemsByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Item, error) {
	query := `
	SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	FROM items
	WHERE order_id = $1;
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("get items by order id: %w", customerr.ErrGetItemsByOrderId)
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var item model.Item
		err = rows.Scan(
			&item.ChrtID, &item.TrackNumber, &item.Price, &item.RID, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("scan item row: %w", customerr.ErrItemScanFailed)
		}

		items = append(items, item)
	}

	return items, nil
}
