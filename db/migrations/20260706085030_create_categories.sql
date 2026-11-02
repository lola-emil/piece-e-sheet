-- +goose Up
SELECT 'up SQL query';
CREATE TABLE categories (
    id UUID PRIMARY KEY,

    user_id UUID NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    name VARCHAR(100) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    deleted_at TIMESTAMPTZ,

    revision BIGINT NOT NULL DEFAULT 1,

    UNIQUE (user_id, name)
);

-- +goose Down
SELECT 'down SQL query';
DROP TABLE categories;