package sql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/creative-rotation/internal/storage/models"
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
	return stats, fmt.Errorf("couldn't get stats: %w", err)
}

func (sr *StatsRepository) StatsSlotSegment(ctx context.Context, slotID, segmentID models.ID) ([]models.Stats, error) {
	var stats []models.Stats
	err := sr.db.SelectContext(ctx, &stats,
		`SELECT * FROM "stats" WHERE slot_id=$1 AND segment_id=$2`, slotID, segmentID)

	return stats, fmt.Errorf("couldn't get stats for slot_id=%d/segment_id=%d: %w", slotID, segmentID, err)
}
