-- +goose Up
SELECT 'up SQL query';
CREATE TABLE expenses (
    id UUID PRIMARY KEY,

    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    category_id UUID
        REFERENCES categories(id)
        ON DELETE SET NULL,

    description VARCHAR(255) NOT NULL,

    amount NUMERIC(12, 2) NOT NULL
        CHECK (amount > 0),

    occurred_at TIMESTAMPTZ NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    deleted_at TIMESTAMPTZ,

    revision BIGINT NOT NULL DEFAULT 1
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE expenses;