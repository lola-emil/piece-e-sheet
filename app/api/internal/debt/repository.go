package debt

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type repo struct{ db *sqlx.DB }
type DebtRepository interface {
	FindAll(ctx context.Context, userID string) ([]Debt, error)
	FindByID(ctx context.Context, id, userID string) (*Debt, error)
	Insert(ctx context.Context, d *Debt) error
	UpdateByID(ctx context.Context, d *Debt) error
	DeleteByID(ctx context.Context, id, userID string) error
}

func NewDebtRepository(db *sqlx.DB) DebtRepository { return &repo{db} }

func (r *repo) FindAll(ctx context.Context, userID string) ([]Debt, error) {
	var debts []Debt
	err := r.db.SelectContext(ctx, &debts,
		`SELECT * FROM debts WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`, userID)
	return debts, err
}

func (r *repo) FindByID(ctx context.Context, id, userID string) (*Debt, error) {
	var d Debt
	err := r.db.GetContext(ctx, &d,
		`SELECT * FROM debts WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`, id, userID)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *repo) Insert(ctx context.Context, d *Debt) error {
	d.ID = uuid.New().String()
	d.RemainingAmount = d.Amount
	d.Status = "pending"
	return r.db.QueryRowxContext(ctx, `
		INSERT INTO debts (id, user_id, person_name, description, amount, remaining_amount, type, due_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING created_at, updated_at, revision`,
		d.ID, d.UserID, d.PersonName, d.Description, d.Amount, d.RemainingAmount, d.Type, d.DueDate).
		Scan(&d.CreatedAt, &d.UpdatedAt, &d.Revision)
}

func (r *repo) UpdateByID(ctx context.Context, d *Debt) error {
	d.UpdatedAt = time.Now()
	rev := d.Revision + 1
	return r.db.QueryRowxContext(ctx, `
		UPDATE debts SET person_name=$1, description=$2, amount=$3, remaining_amount=$4, 
		type=$5, status=$6, due_date=$7, settled_at=$8, updated_at=NOW(), revision=$9
		WHERE id=$10 AND user_id=$11 AND deleted_at IS NULL AND revision=$12
		RETURNING updated_at, revision`,
		d.PersonName, d.Description, d.Amount, d.RemainingAmount, d.Type, d.Status, d.DueDate, d.SettledAt, rev, d.ID, d.UserID, d.Revision).
		Scan(&d.UpdatedAt, &d.Revision)
}

func (r *repo) DeleteByID(ctx context.Context, id, userID string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE debts SET deleted_at=NOW() WHERE id=$1 AND user_id=$2 AND deleted_at IS NULL`, id, userID)
	if err != nil {
		return err
	}
	if rows, _ := res.RowsAffected(); rows == 0 {
		return sql.ErrNoRows
	}
	return nil
}
