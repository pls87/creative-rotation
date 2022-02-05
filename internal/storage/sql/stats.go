package sql

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type statsKind int

const (
	impressions statsKind = iota
	conversions
)

func (d statsKind) String() string {
	return [...]string{"impressions", "conversions"}[d]
}

type StatsRepository struct {
	db *sqlx.DB
}

func (sr *StatsRepository) StatsSlotSegment(ctx context.Context, slotID, segmentID models.ID) ([]models.Stats, error) {
	var stats []models.Stats
	err := sr.db.SelectContext(ctx, &stats,
		`SELECT * FROM "stats" WHERE slot_id=$1 AND segment_id=$2`, slotID, segmentID)

	return stats, fmt.Errorf("couldn't get stats for slot_id=%d/segment_id=%d: %w", slotID, segmentID, err)
}

func (sr *StatsRepository) UpdateStatsImpression(ctx context.Context, impression models.Impression) error {
	return sr.updateStats(ctx, impressions, impression.CreativeID, impression.SlotID, impression.SegmentID)
}

func (sr *StatsRepository) UpdateStatsConversion(ctx context.Context, conversion models.Conversion) error {
	return sr.updateStats(ctx, conversions, conversion.CreativeID, conversion.SlotID, conversion.SegmentID)
}

func (sr *StatsRepository) updateStats(ctx context.Context, kind statsKind,
	creativeID, slotID, segmentID models.ID) error {
	query := `INSERT INTO "stats" ("$1", creative_id, slot_id, segment_id) 
		VALUES (1, $2, $3, $4)
		ON CONFLICT (creative_id, slot_id, segment_id) DO UPDATE SET "$1" = "$1" + 1;`
	_, err := sr.db.ExecContext(ctx, query, kind.String(), creativeID, slotID, segmentID)
	if err != nil {
		return fmt.Errorf("couldn't update stats for creative_id=%d, slot_id=%d, segment_id=%d: %w",
			creativeID, slotID, segmentID, err)
	}

	return nil
}
