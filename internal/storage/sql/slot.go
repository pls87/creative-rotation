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
	if err := sr.db.SelectContext(ctx, &slots, ALLQuery("slot")); err != nil {
		return nil, ALLError("slot", err)
	}

	return slots, nil
}

func (sr *SlotRepository) Create(ctx context.Context, s models.Slot) (added models.Slot, err error) {
	var lastInsertID int
	if err = sr.db.QueryRowxContext(ctx, CREATEQuery("slot"), s.Desc).Scan(&lastInsertID); err != nil {
		return s, CREATEError("slot", err)
	}

	s.ID = models.ID(lastInsertID)
	return s, nil
}

func (sr *SlotRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := sr.db.ExecContext(ctx, DELETEQuery("slot"), id)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return DELETEError("slot", id, basic.ErrDoesNotExist)
		}
		return nil
	}
	return DELETEError("slot", id, err)
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
