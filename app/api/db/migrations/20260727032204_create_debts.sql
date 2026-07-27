-- +goose Up
SELECT 'up SQL query';
CREATE TABLE debts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    
    person_name VARCHAR(100) NOT NULL,
    description VARCHAR(255),
    
    amount NUMERIC(12, 2) NOT NULL CHECK (amount > 0),
    remaining_amount NUMERIC(12, 2) NOT NULL CHECK (remaining_amount >= 0),
    
    -- 'owed_to_me' (they owe me) or 'i_owe' (I owe them)
    type VARCHAR(10) NOT NULL CHECK (type IN ('owed_to_me', 'i_owe')),
    
    status VARCHAR(20) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'partial', 'paid')),
    
    due_date TIMESTAMPTZ,
    settled_at TIMESTAMPTZ,
    
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMPTZ,
    revision BIGINT NOT NULL DEFAULT 1
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE debts;