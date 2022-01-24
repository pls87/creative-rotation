package sql

import (
	"context"
	"fmt"

	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"

	"github.com/jmoiron/sqlx"
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

	return slots, err
}

func (sr *SlotRepository) Create(ctx context.Context, s models.Slot) (added models.Slot, err error) {
	query := `INSERT INTO "slot" (description) VALUES ('?')`
	res, err := sr.db.ExecContext(ctx, query, s.Desc)
	if err == nil {
		id, _ := res.LastInsertId()
		s.ID = models.ID(id)
	}

	return s, err
}

func (sr *SlotRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := sr.db.ExecContext(ctx, `DELETE FROM "slot" WHERE ID=?`, id)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return fmt.Errorf("DELETE: slot id=%d: %w", id, basic.ErrDoesNotExist)
		}
	}
	return err
}