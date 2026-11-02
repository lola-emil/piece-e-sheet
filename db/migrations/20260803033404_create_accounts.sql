-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    type VARCHAR(50) NOT NULL DEFAULT 'checking', -- e.g., cash, checking, credit_card, savings
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    revision BIGINT NOT NULL DEFAULT 1,
    UNIQUE (user_id, name)
);

ALTER TABLE expenses ADD COLUMN account_id UUID REFERENCES accounts(id) ON DELETE SET NULL;

-- +goose StatementEnd

-- +goose Down
SELECT 'down SQL query';
-- +goose StatementBegin
ALTER TABLE expenses
    DROP COLUMN IF EXISTS account_id;

DROP TABLE IF EXISTS accounts;
-- +goose StatementEnd
