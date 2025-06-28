package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/aliskhannn/order-service/internal/model"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) SaveOrder(ctx context.Context, order *model.Order) (string, error) {
	tx, err := r.db.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", fmt.Errorf("error beginning tx: %w", err)
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
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
		return "", fmt.Errorf("error inserting into orders: %w", err)
	}

	// Insert into delivery
	d := order.Delivery
	deleviryQuery := `
	INSERT INTO delivery (
	    order_uid, name, phone, zip, city, address, region, email
	) VALUES($1, $2, $3, $4, $5, $6, $7, $8);
	`

	_, err = tx.Exec(ctx, deleviryQuery,
		order.OrderID, d.Name, d.Phone, d.Zip, d.City, d.Address, d.Region, d.Email)
	if err != nil {
		return "", fmt.Errorf("error inserting into delivery: %w", err)
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
		return "", fmt.Errorf("error inserting into deliveryt: %w", err)
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
			return "", fmt.Errorf("error inserting into items: %w", err)
		}
	}

	return order.OrderID, nil
}

func (r *OrderRepository) GetOrderById(ctx context.Context, orderID string) (*model.Order, error) {
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
			return nil, fmt.Errorf("no order found with id %s", orderID)
		}
		return nil, fmt.Errorf("error getting order by id %s: %v", orderID, err)
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
			return nil, fmt.Errorf("no delivery found for order %s", orderID)
		}
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
			return nil, fmt.Errorf("no payment found for order %s", orderID)
		}
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
			return nil, fmt.Errorf("error scanning item for order %s: %v", orderID, err)
		}
		items = append(items, item)
	}
	order.Items = items

	return &order, nil
}
