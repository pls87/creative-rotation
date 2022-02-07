package sql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type SegmentRepository struct {
	db *sqlx.DB
}

func (sr *SegmentRepository) All(ctx context.Context) ([]models.Segment, error) {
	var segments []models.Segment
	if err := sr.db.SelectContext(ctx, &segments, `SELECT * FROM "segment"`); err != nil {
		return nil, fmt.Errorf("couldn't get segments from database: %w", err)
	}

	return segments, nil
}

func (sr *SegmentRepository) Create(ctx context.Context, s models.Segment) (added models.Segment, err error) {
	query := `INSERT INTO "segment" (description) VALUES ($1) RETURNING "ID"`
	lastInsertID := 0
	if err = sr.db.QueryRowxContext(ctx, query, s.Desc).Scan(&lastInsertID); err != nil {
		return s, fmt.Errorf("couldn't create segment in database: %w", err)
	}

	s.ID = models.ID(lastInsertID)
	return s, nil
}

func (sr *SegmentRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := sr.db.ExecContext(ctx, `DELETE FROM "segment" WHERE "ID"=$1`, id)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return fmt.Errorf("couldn't delete segment id=%d: %w", id, basic.ErrDoesNotExist)
		}
		return nil
	}
	return fmt.Errorf("couldn't delete segment id=%d: %w", id, err)
}
