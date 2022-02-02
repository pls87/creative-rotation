//nolint: dupl
package sql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type SlotRepository struct {
	db *sqlx.DB
}

func (sr *SlotRepository) Init(ctx context.Context) error {
	return nil
}

func (sr *SlotRepository) All(ctx context.Context) ([]models.Slot, error) {
	var slots []models.Slot
	err := sr.db.SelectContext(ctx, &slots, `SELECT * FROM "slot"`)

	return slots, fmt.Errorf("couldn't get slots from database: %w", err)
}

func (sr *SlotRepository) Create(ctx context.Context, s models.Slot) (added models.Slot, err error) {
	query := `INSERT INTO "slot" (description) VALUES ($1) RETURNING "ID"`
	lastInsertID := 0
	err = sr.db.QueryRowxContext(ctx, query, s.Desc).Scan(&lastInsertID)
	if err == nil {
		s.ID = models.ID(lastInsertID)
	}
	return s, fmt.Errorf("couldn't create slot in database: %w", err)
}

func (sr *SlotRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := sr.db.ExecContext(ctx, `DELETE FROM "slot" WHERE "ID"=$1`, id)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return fmt.Errorf("couldn't delete slot id=%d: %w", id, basic.ErrDoesNotExist)
		}
	}
	return fmt.Errorf("couldn't delete slot id=%d: %w", id, err)
}
