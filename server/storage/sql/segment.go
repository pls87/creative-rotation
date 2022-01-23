package sql

import (
	"context"
	"fmt"

	"github.com/pls87/creative-rotation/server/storage/basic"
	"github.com/pls87/creative-rotation/server/storage/models"

	"github.com/jmoiron/sqlx"
)

type SegmentRepository struct {
	db *sqlx.DB
}

func (sr *SegmentRepository) Init(ctx context.Context) error {
	return nil
}

func (sr *SegmentRepository) All(ctx context.Context) ([]models.Segment, error) {
	var segments []models.Segment
	err := sr.db.SelectContext(ctx, &segments, `SELECT * FROM "segment"`)

	return segments, err
}

func (sr *SegmentRepository) Create(ctx context.Context, s models.Segment) (added models.Segment, err error) {
	query := `INSERT INTO "segment" (description) VALUES ('?')`
	res, err := sr.db.ExecContext(ctx, query, s.Desc)
	if err == nil {
		id, _ := res.LastInsertId()
		s.ID = models.ID(id)
	}

	return s, err
}

func (sr *SegmentRepository) Delete(ctx context.Context, id models.ID) error {
	res, err := sr.db.ExecContext(ctx, `DELETE FROM "segment" WHERE ID=?`, id)
	if err == nil {
		if affected, _ := res.RowsAffected(); affected == 0 {
			return fmt.Errorf("DELETE: segment id=%d: %w", id, basic.ErrDoesNotExist)
		}
	}
	return err
}
