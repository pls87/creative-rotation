package sql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
	"github.com/pls87/creative-rotation/internal/storage/sql/errors"
	"github.com/pls87/creative-rotation/internal/storage/sql/queries"
)

type SegmentRepository struct {
	db *sqlx.DB
}

func (sr *SegmentRepository) All(ctx context.Context) ([]models.Segment, error) {
	var segments []models.Segment
	if err := sr.db.SelectContext(ctx, &segments, queries.Crud.All(queries.SegmentRelation)); err != nil {
		return nil, errors.Crud.All(queries.SegmentRelation, err)
	}

	return segments, nil
}

func (sr *SegmentRepository) Create(ctx context.Context, s models.Segment) (added models.Segment, err error) {
	var lastInsertID int
	if err = sr.db.QueryRowxContext(ctx, queries.Crud.Create(queries.SegmentRelation), s.Desc).
		Scan(&lastInsertID); err != nil {
		return s, errors.Crud.Create(queries.SegmentRelation, err)
	}

	s.ID = models.ID(lastInsertID)
	return s, nil
}

func (sr *SegmentRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := sr.db.ExecContext(ctx, queries.Crud.Delete(queries.SegmentRelation), id)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return errors.Crud.Delete(queries.SegmentRelation, id, basic.ErrDoesNotExist)
		}
		return nil
	}
	return errors.Crud.Delete(queries.SegmentRelation, id, err)
}
