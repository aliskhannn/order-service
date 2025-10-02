package order

import (
	"context"
	"errors"
	"fmt"

	"go.uber.org/zap"

	"github.com/google/uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aliskhannn/order-service/internal/model"
)

// Repository provides access to order-related data stored in PostgreSQL.
// It encapsulates CRUD operations for orders, deliveries, payments, and items.
type Repository struct {
	db     *pgxpool.Pool
	logger *zap.Logger
}

// New creates a new Repository instance with a given pgx connection pool and logger.
func New(db *pgxpool.Pool, lg *zap.Logger) *Repository {
	return &Repository{
		db:     db,
		logger: lg,
	}
}

// SaveOrder inserts a new order and all its related entities (delivery, payment, items)
// into the database in a single transaction. If any step fails, the transaction is rolled back.
// Returns the generated order UUID or an error.
func (r *Repository) SaveOrder(ctx context.Context, order *model.Order) (uuid.UUID, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed begin tx: %w", err)
	}
	// Ensure the transaction will be rolled back if not committed
	defer func() {
		_ = tx.Rollback(ctx)
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
		return uuid.Nil, fmt.Errorf("failed to create order: %w", err)
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
		return uuid.Nil, fmt.Errorf("failed to insert delivery: %w", err)
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
		return uuid.Nil, fmt.Errorf("failed to insert payment: %w", err)
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
			return uuid.Nil, fmt.Errorf("failed to insert item: %w", err)
		}
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return uuid.Nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order.OrderID, nil
}

// GetOrderById retrieves a single order by its UUID, including delivery and payment details.
// Note: Items are not currently loaded here — use GetItemsByOrderID if needed.
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
			return nil, fmt.Errorf("failed to get order by id: %w", ErrOrderNotFound)
		}

		return nil, fmt.Errorf("failed to scan row: %w", err)
	}

	o.Delivery = d
	o.Payment = p

	return &o, err
}

// GetLastOrders fetches the most recent orders from the database, limited by the given number.
// Each order includes its delivery, payment, and items.
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
		return nil, fmt.Errorf("failed to get last orders: %w", err)
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
			return nil, fmt.Errorf("failed to scan order row: %w", err)
		}

		o.Delivery = d
		o.Payment = p

		items, err := r.GetItemsByOrderID(ctx, o.OrderID)
		if err != nil {
			return nil, err
		}

		o.Items = items

		orders = append(orders, o)
	}

	if len(orders) == 0 {
		return []model.Order{}, nil
	}

	return orders, nil
}

// GetItemsByOrderID retrieves all items associated with the given order UUID.
// Returns an empty slice if no items are found.
func (r *Repository) GetItemsByOrderID(ctx context.Context, orderID uuid.UUID) ([]model.Item, error) {
	query := `
	SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status
	FROM items
	WHERE order_id = $1;
	`

	rows, err := r.db.Query(ctx, query, orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get items by order id: %w", err)
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
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		items = append(items, item)
	}

	if len(items) == 0 {
		return []model.Item{}, nil
	}

	return items, nil
}
