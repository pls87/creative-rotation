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

func (sr *SlotRepository) All(ctx context.Context) ([]models.Slot, error) {
	var slots []models.Slot
	if err := sr.db.SelectContext(ctx, &slots, `SELECT * FROM "slot"`); err != nil {
		return nil, fmt.Errorf("couldn't get slots from database: %w", err)
	}

	return slots, nil
}

func (sr *SlotRepository) Create(ctx context.Context, s models.Slot) (added models.Slot, err error) {
	query := `INSERT INTO "slot" (description) VALUES ($1) RETURNING "ID"`
	lastInsertID := 0
	if err = sr.db.QueryRowxContext(ctx, query, s.Desc).Scan(&lastInsertID); err != nil {
		return s, fmt.Errorf("couldn't create slot in database: %w", err)
	}

	s.ID = models.ID(lastInsertID)
	return s, nil
}

func (sr *SlotRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := sr.db.ExecContext(ctx, `DELETE FROM "slot" WHERE "ID"=$1`, id)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return fmt.Errorf("couldn't delete slot id=%d: %w", id, basic.ErrDoesNotExist)
		}
		return nil
	}
	return fmt.Errorf("couldn't delete slot id=%d: %w", id, err)
}

func (sr *SlotRepository) Creatives(ctx context.Context, id models.ID) ([]models.Creative, error) {
	var creatives []models.Creative
	query := `SELECT s.* FROM "slot_creative" sc
		INNER JOIN "creative" cr ON sc.creative_id=cr."ID" 
		WHERE sc.slot_id = $1`
	if err := sr.db.SelectContext(ctx, &creatives, query, id); err != nil {
		return nil, fmt.Errorf("couldn't get creatives for slot '%d' from database: %w", id, err)
	}

	return creatives, nil
}
