-- +goose Up
SELECT 'up SQL query';
CREATE TABLE payments (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    debt_id UUID NOT NULL REFERENCES debts(id) ON DELETE CASCADE,
    
    amount NUMERIC(12, 2) NOT NULL CHECK (amount > 0),
    note VARCHAR(255),
    paid_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    revision BIGINT NOT NULL DEFAULT 1
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE payments;