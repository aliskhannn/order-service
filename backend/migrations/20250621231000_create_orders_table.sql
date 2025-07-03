-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS orders (
    order_uid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    track_number VARCHAR(32) NOT NULL,
    entry VARCHAR(10),
    locale VARCHAR(5),
    internal_signature TEXT,
    customer_id VARCHAR(50),
    delivery_service VARCHAR(50),
    shardkey VARCHAR(5),
    sm_id INT,
    date_created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    oof_shard VARCHAR(5)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP EXTENSION IF EXISTS "uuid-ossp";
DROP TABLE IF EXISTS orders;
-- +goose StatementEnd
