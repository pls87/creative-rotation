package sql

import (
	"context"

	"github.com/pls87/creative-rotation/internal/storage/models"

	"github.com/jmoiron/sqlx"
)

type StatsRepository struct {
	db *sqlx.DB
}

func (sr *StatsRepository) Init(ctx context.Context) error {
	return nil
}

func (sr *StatsRepository) AllStats(ctx context.Context) ([]models.Stats, error) {
	var stats []models.Stats
	err := sr.db.SelectContext(ctx, &stats, `SELECT * FROM "stats"`)
	return stats, err
}

func (sr *StatsRepository) StatsSlotSegment(ctx context.Context, slotId, segmentId models.ID) ([]models.Stats, error) {
	var stats []models.Stats
	err := sr.db.SelectContext(ctx, &stats,
		`SELECT * FROM "stats" WHERE slot_id=? AND segment_id=?`, slotId, segmentId)
	return stats, err
}
