package sql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type CreativeRepository struct {
	db *sqlx.DB
}

func (cr *CreativeRepository) All(ctx context.Context) ([]models.Creative, error) {
	var creatives []models.Creative
	if err := cr.db.SelectContext(ctx, &creatives, ALLQuery("creative")); err != nil {
		return nil, ALLError("creative", err)
	}

	return creatives, nil
}

func (cr *CreativeRepository) Create(ctx context.Context, c models.Creative) (added models.Creative, err error) {
	var lastInsertID int
	if err = cr.db.QueryRowxContext(ctx, CREATEQuery("creative"), c.Desc).Scan(&lastInsertID); err != nil {
		return c, CREATEError("creative", err)
	}

	c.ID = models.ID(lastInsertID)
	return c, nil
}

func (cr *CreativeRepository) Slots(ctx context.Context, id models.ID) ([]models.Slot, error) {
	var slots []models.Slot
	query := `SELECT s.* FROM "slot_creative" sc
		INNER JOIN "slot" s ON sc.slot_id=s."ID"
		WHERE sc.creative_id = $1`
	if err := cr.db.SelectContext(ctx, &slots, query, id); err != nil {
		return nil, fmt.Errorf("couldn't get slots for creative '%d' from database: %w", id, err)
	}

	return slots, nil
}

func (cr *CreativeRepository) AllCreativeSlots(ctx context.Context) ([]models.SlotCreative, error) {
	var slotCreatives []models.SlotCreative
	query := `SELECT sc.slot_id, sc.creative_id, s.description as slot_desc, cr.description as creative_desc 
		FROM "slot_creative" sc 
		INNER JOIN "slot" s ON sc.slot_id=s."ID"
		INNER JOIN "creative" cr ON sc.creative_id=cr."ID"`
	if err := cr.db.SelectContext(ctx, &slotCreatives, query); err != nil {
		return nil, fmt.Errorf("couldn't get slots-creative from database: %w", err)
	}

	return slotCreatives, nil
}

func (cr *CreativeRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := cr.db.ExecContext(ctx, DELETEQuery("creative"), id)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return DELETEError("creative", id, basic.ErrDoesNotExist)
		}
		return nil
	}
	return DELETEError("creative", id, err)
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

	return rows.Next(), nil
}
