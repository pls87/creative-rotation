package sql

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/creative-rotation/internal/storage/basic"
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

	return creatives, fmt.Errorf("couldn't get creatives from database: %w", err)
}

func (cr *CreativeRepository) Create(ctx context.Context, c models.Creative) (added models.Creative, err error) {
	query := `INSERT INTO "creative" (description) VALUES ($1) RETURNING "ID"`
	lastInsertID := 0
	err = cr.db.QueryRowxContext(ctx, query, c.Desc).Scan(&lastInsertID)
	if err == nil {
		c.ID = models.ID(lastInsertID)
	}

	return c, fmt.Errorf("couldn't create creative in database: %w", err)
}

func (cr *CreativeRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := cr.db.ExecContext(ctx, `DELETE FROM "creative" WHERE "ID"=$1`, id)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return fmt.Errorf("couldn't delete creative id=%d: %w", id, basic.ErrDoesNotExist)
		}
	}
	return fmt.Errorf("couldn't delete creative id=%d: %w", id, err)
}

func (cr *CreativeRepository) ToSlot(ctx context.Context, creativeID, slotID models.ID) error {
	query := `INSERT INTO "slot_creative" (creative_id, slot_id) VALUES ($1, $2)`
	res, err := cr.db.ExecContext(ctx, query, creativeID, slotID)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("couldn't add creative id=%d to slot_id=%d: %w",
			creativeID, slotID, basic.ErrCreativeAlreadyInSlot)
	}

	return nil
}

func (cr *CreativeRepository) FromSlot(ctx context.Context, creativeID, slotID models.ID) error {
	res, err := cr.db.ExecContext(ctx, `DELETE FROM "slot_creative" WHERE creative_id = $1 AND slot_id=$2`,
		creativeID, slotID)
	if err != nil {
		return err
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return fmt.Errorf("couldn't delete creative id=%d from slot_id=%d: %w",
			creativeID, slotID, basic.ErrCreativeNotInSlot)
	}

	return nil
}

func (cr *CreativeRepository) InSlot(ctx context.Context, creativeID, slotID models.ID) (bool, error) {
	rows, err := cr.db.QueryxContext(ctx, `SELECT * FROM "slot_creative" WHERE creative_id = $1 AND slot_id = $2`,
		creativeID, slotID)
	if err != nil {
		return false, fmt.Errorf("couldn't get info about slot/creative: %w", err)
	}
	defer rows.Close()

	return rows.Next(), fmt.Errorf("couldn't get info about slot/creative: %w", err)
}

func (cr *CreativeRepository) TrackImpression(ctx context.Context, imp models.Impression) error {
	query := `INSERT INTO "impression" (creative_id, slot_id, segment_id, time) VALUES ($1, $2, $3, $4)`
	_, err := cr.db.ExecContext(ctx, query, imp.CreativeID, imp.SlotID, imp.SegmentID, time.Now())

	return err
}

func (cr *CreativeRepository) TrackConversion(ctx context.Context, conversion models.Conversion) error {
	query := `INSERT INTO "conversion" (creative_id, slot_id, segment_id, time) VALUES ($1, $2, $3, $4)`
	_, err := cr.db.ExecContext(ctx, query, conversion.CreativeID, conversion.SlotID, conversion.SegmentID, time.Now())

	return err
}
