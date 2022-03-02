package sql

import (
	"context"

	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
	"github.com/pls87/creative-rotation/internal/storage/sql/errors"
	"github.com/pls87/creative-rotation/internal/storage/sql/queries"
)

type SlotRepository struct {
	db DB
}

func (sr *SlotRepository) All(ctx context.Context) ([]models.Slot, error) {
	var slots []models.Slot
	if err := sr.db.SelectContext(ctx, &slots, queries.Crud.All(queries.SlotRelation)); err != nil {
		return nil, errors.Crud.All("slot", err)
	}

	return slots, nil
}

func (sr *SlotRepository) Create(ctx context.Context, s models.Slot) (models.Slot, error) {
	id, err := sr.db.InsertRow(ctx, queries.Crud.Create(queries.SlotRelation), s.Desc)
	if err != nil {
		return s, errors.Crud.Create(queries.SlotRelation, err)
	}

	s.ID = id
	return s, nil
}

func (sr *SlotRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := sr.db.ExecContext(ctx, queries.Crud.Delete(queries.SlotRelation), id)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return errors.Crud.Delete(queries.SlotRelation, id, basic.ErrDoesNotExist)
		}
		return nil
	}
	return errors.Crud.Delete(queries.SlotRelation, id, err)
}

func (sr *SlotRepository) Creatives(ctx context.Context, id models.ID) ([]models.Creative, error) {
	var creatives []models.Creative
	query := queries.SC.GetFor(queries.SlotRelation, queries.CreativeRelation)
	if err := sr.db.SelectContext(ctx, &creatives, query, id); err != nil {
		return nil, errors.SC.GetFor(queries.SlotRelation, queries.CreativeRelation, id, err)
	}

	return creatives, nil
}
