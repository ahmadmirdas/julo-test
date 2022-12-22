-- +goose Up
-- +goose StatementBegin
CREATE TABLE wallet
(
    id uuid DEFAULT gen_random_uuid (),
    owned_by uuid NOT NULL UNIQUE,
    is_enabled BOOLEAN NOT NULL DEFAULT true,
    balance FLOAT NOT NULL DEFAULT 0,
    enabled_at TIMESTAMP NULL,
    disabled_at TIMESTAMP NULL,

    PRIMARY KEY (id)
);

CREATE INDEX idx_wallet_owned_by ON wallet(owned_by);

CREATE TABLE history
(
    id uuid DEFAULT gen_random_uuid (),
    wallet_id uuid NOT NULL,
    status VARCHAR DEFAULT 'pending',
    type VARCHAR NOT NULL,
    amount FLOAT NOT NULL DEFAULT 0,
    reference_id uuid NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    PRIMARY KEY (id),
    CONSTRAINT fk_history_wallet_id FOREIGN KEY (wallet_id) REFERENCES "wallet" (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE history;
DROP TABLE wallet;
-- +goose StatementEnd