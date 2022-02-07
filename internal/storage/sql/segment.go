package sql

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type SegmentRepository struct {
	db *sqlx.DB
}

func (sr *SegmentRepository) All(ctx context.Context) ([]models.Segment, error) {
	var segments []models.Segment
	if err := sr.db.SelectContext(ctx, &segments, ALLQuery("segment")); err != nil {
		return nil, ALLError("segment", err)
	}

	return segments, nil
}

func (sr *SegmentRepository) Create(ctx context.Context, s models.Segment) (added models.Segment, err error) {
	var lastInsertID int
	if err = sr.db.QueryRowxContext(ctx, CREATEQuery("segment"), s.Desc).Scan(&lastInsertID); err != nil {
		return s, CREATEError("segment", err)
	}

	s.ID = models.ID(lastInsertID)
	return s, nil
}

func (sr *SegmentRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := sr.db.ExecContext(ctx, DELETEQuery("segment"), id)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return DELETEError("segment", id, basic.ErrDoesNotExist)
		}
		return nil
	}
	return DELETEError("segment", id, err)
}
