package repository

import (
	"context"
	"fmt"
	"github.com/aliskhannn/order-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order *model.Order) (uuid.UUID, error)
	GetOrderById(ctx context.Context, orderID uuid.UUID) (*model.Order, error)
}

type orderRepo struct {
	db *pgxpool.Pool
}

func NewOrderRepo(db *pgxpool.Pool) OrderRepository {
	return &orderRepo{db: db}
}

func (r *orderRepo) SaveOrder(ctx context.Context, order *model.Order) (uuid.UUID, error) {
	query := `
	INSERT INTO orders (
		order_uid, track_number, entry, locale, internal_signature, customer_id,
		delivery_service, shardkey, sm_id, date_created, oof_shard
	) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING order_uid;
	`

	var orderID uuid.UUID
	err := r.db.QueryRow(ctx, query,
		order.OrderUid, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard).Scan(&orderID)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error creatin order: %w", err)
	}

	return orderID, nil
}

func (r *orderRepo) GetOrderById(ctx context.Context, orderID uuid.UUID) (*model.Order, error) {
	query := `
	SELECT order_uid, track_number, entry, locale, internal_signature, customer_id,
		delivery_service, shardkey, sm_id, date_created, oof_shard
	FROM orders
	WHERE order_uid = $1;
	`

	var order model.Order
	err := r.db.QueryRow(ctx, query, orderID).Scan(
		&order.OrderUid, &order.TrackNumber, &order.Entry, &order.Locale,
		&order.InternalSignature, &order.CustomerId, &order.DeliveryService,
		&order.Shardkey, &order.SmId, &order.DateCreated, &order.OofShard,
	)
	if err != nil {
		return nil, fmt.Errorf("error getting order by id %d: %v", orderID, err)
	}

	return &order, nil
}
