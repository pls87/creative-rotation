package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/pls87/creative-rotation/internal/storage/basic"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type CreativeRepository struct {
	db *sqlx.DB
}

func (cr *CreativeRepository) Init(ctx context.Context) error {
	return nil
}

func (cr *CreativeRepository) All(ctx context.Context) ([]models.Creative, error) {
	var creatives []models.Creative
	err := cr.db.SelectContext(ctx, &creatives, `SELECT * FROM "creative"`)

	return creatives, err
}

func (cr *CreativeRepository) Create(ctx context.Context, c models.Creative) (added models.Creative, err error) {
	query := `INSERT INTO "creative" (description) VALUES ('?')`
	res, err := cr.db.ExecContext(ctx, query, c.Desc)
	if err == nil {
		id, _ := res.LastInsertId()
		c.ID = models.ID(id)
	}

	return c, err
}

func (cr *CreativeRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := cr.db.ExecContext(ctx, `DELETE FROM "creative" WHERE ID=?`, id)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return fmt.Errorf("DELETE: creative id=%d: %w", id, basic.ErrDoesNotExist)
		}
	}
	return err
}

func (cr *CreativeRepository) ToSlot(ctx context.Context, creativeId, slotId models.ID) error {
	query := `INSERT INTO "slot_creative" (creative_id, slot_id) VALUES (?, ?)`
	res, err := cr.db.ExecContext(ctx, query, creativeId, slotId)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("adding to slot: creative id=%d already in slot_id=%d: %w",
			creativeId, slotId, basic.ErrCreativeAlreadyInSlot)
	}

	return nil
}

func (cr *CreativeRepository) FromSlot(ctx context.Context, creativeId, slotId models.ID) error {
	res, err := cr.db.ExecContext(ctx, `DELETE FROM "slot_creative" WHERE creative_id = ? AND slot_id=?`,
		creativeId, slotId)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("removing from slot: creative id=%d not in slot_id=%d: %w",
			creativeId, slotId, basic.ErrCreativeNotInSlot)
	}

	return nil
}

func (cr *CreativeRepository) InSlot(ctx context.Context, creativeId, slotId models.ID) (bool, error) {
	rows, err := cr.db.QueryxContext(ctx, `SELECT * FROM "slot_creative" WHERE creative_id = ? AND slot_id = ?`,
		creativeId, slotId)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	return rows.Next(), err
}

func (cr *CreativeRepository) TrackImpression(ctx context.Context, creativeId, slotId, segmentId models.ID) error {
	query := `INSERT INTO "impression" (creative_id, slot_id, segment_id, time) VALUES (?, ?, ?, ?)`
	_, err := cr.db.ExecContext(ctx, query, creativeId, slotId, segmentId, time.Now())

	return err
}

func (cr *CreativeRepository) TrackConversion(ctx context.Context, creativeId, slotId, segmentId models.ID) error {
	query := `INSERT INTO "conversion" (creative_id, slot_id, segment_id, time) VALUES (?, ?, ?, ?)`
	_, err := cr.db.ExecContext(ctx, query, creativeId, slotId, segmentId, time.Now())

	return err
}
