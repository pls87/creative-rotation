package sql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
	"github.com/pls87/creative-rotation/internal/storage/sql/errors"
	"github.com/pls87/creative-rotation/internal/storage/sql/queries"
)

type CreativeRepository struct {
	db *sqlx.DB
}

func (cr *CreativeRepository) All(ctx context.Context) ([]models.Creative, error) {
	var creatives []models.Creative
	if err := cr.db.SelectContext(ctx, &creatives, queries.Crud.All(queries.CreativeRelation)); err != nil {
		return nil, errors.Crud.All(queries.CreativeRelation, err)
	}

	return creatives, nil
}

func (cr *CreativeRepository) Create(ctx context.Context, c models.Creative) (added models.Creative, err error) {
	var lastInsertID int
	if err = cr.db.QueryRowxContext(ctx, queries.Crud.Create(queries.CreativeRelation), c.Desc).
		Scan(&lastInsertID); err != nil {
		return c, errors.Crud.Create(queries.CreativeRelation, err)
	}

	c.ID = models.ID(lastInsertID)
	return c, nil
}

func (cr *CreativeRepository) Slots(ctx context.Context, id models.ID) ([]models.Slot, error) {
	var slots []models.Slot
	query := queries.SC.GetFor(queries.CreativeRelation, queries.SlotRelation)
	if err := cr.db.SelectContext(ctx, &slots, query, id); err != nil {
		return nil, errors.SC.GetFor(queries.CreativeRelation, queries.SlotRelation, id, err)
	}

	return slots, nil
}

func (cr *CreativeRepository) AllCreativeSlots(ctx context.Context) ([]models.SlotCreative, error) {
	var slotCreatives []models.SlotCreative
	if err := cr.db.SelectContext(ctx, &slotCreatives, queries.SC.All()); err != nil {
		return nil, errors.Crud.All(queries.LocationRelation, err)
	}

	return slotCreatives, nil
}

func (cr *CreativeRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := cr.db.ExecContext(ctx, queries.Crud.Delete(queries.CreativeRelation), id)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return errors.Crud.Delete(queries.CreativeRelation, id, basic.ErrDoesNotExist)
		}
		return nil
	}
	return errors.Crud.Delete(queries.CreativeRelation, id, err)
}

func (cr *CreativeRepository) ToSlot(ctx context.Context, creativeID, slotID models.ID) error {
	res, err := cr.db.ExecContext(ctx, queries.SC.Create(), creativeID, slotID)
	if err != nil {
		return errors.SC.Create(creativeID, slotID, err)
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return errors.SC.Create(creativeID, slotID, basic.ErrCreativeAlreadyInSlot)
	}

	return nil
}

func (cr *CreativeRepository) FromSlot(ctx context.Context, creativeID, slotID models.ID) error {
	res, err := cr.db.ExecContext(ctx, queries.SC.Delete(), creativeID, slotID)
	if err != nil {
		return errors.SC.Delete(creativeID, slotID, err)
	}
	if affected, _ := res.RowsAffected(); affected == 0 {
		return errors.SC.Delete(creativeID, slotID, basic.ErrCreativeNotInSlot)
	}

	return nil
}

func (cr *CreativeRepository) InSlot(ctx context.Context, creativeID, slotID models.ID) (bool, error) {
	rows, err := cr.db.QueryxContext(ctx, queries.SC.Exists(), creativeID, slotID)
	if err != nil {
		return false, errors.SC.Exists(creativeID, slotID, err)
	}

	defer rows.Close()

	return rows.Next(), nil
}
