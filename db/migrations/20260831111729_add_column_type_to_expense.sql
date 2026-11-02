-- +goose Up
-- +goose StatementBegin
SELECT
    'up SQL query';

ALTER TABLE
    expenses
ADD
    COLUMN type VARCHAR(10) NOT NULL DEFAULT 'expense';

ALTER TABLE
    expenses
ADD
    CONSTRAINT check_expense_type CHECK (type IN ('expense', 'income'));

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
SELECT
    'down SQL query';
ALTER TABLE expenses DROP COLUMN type;
-- +goose StatementEnd