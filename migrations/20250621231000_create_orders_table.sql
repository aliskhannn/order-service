-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS orders (
    order_uid TEXT PRIMARY KEY,
    track_number VARCHAR(32) NOT NULL,
    entry VARCHAR(10),
    locale VARCHAR(5),
    internal_signature TEXT,
    customer_id VARCHAR(50),
    delivery_service VARCHAR(50),
    shardkey VARCHAR(5),
    sm_id INT,
    date_created TIMESTAMPTZ NOT NULL,
    oof_shard VARCHAR(5)
);

-- +goose StatementEnd
-- {"order_id": 1, "status": "created"}
-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
