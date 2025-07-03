-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS items (
    id SERIAL PRIMARY KEY,
    order_id UUID NOT NULL REFERENCES orders(order_uid) ON DELETE CASCADE,
    chrt_id BIGINT NOT NULL,
    track_number VARCHAR(32) NOT NULL,
    price INT NOT NULL,
    rid TEXT NOT NULL,
    name VARCHAR(100) NOT NULL,
    sale INT NOT NULL,
    size VARCHAR(10),
    total_price INT NOT NULL,
    nm_id BIGINT NOT NULL,
    brand VARCHAR(100) NOT NULL,
    status INT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS items
-- +goose StatementEnd
