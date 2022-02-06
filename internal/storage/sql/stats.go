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

func (k statsKind) String() string {
	return [...]string{"impressions", "conversions"}[k]
}

type StatsRepository struct {
	db *sqlx.DB
}

func (sr *StatsRepository) StatsSlotSegment(ctx context.Context, slotID, segmentID models.ID) ([]models.Stats, error) {
	var stats []models.Stats
	err := sr.db.SelectContext(ctx, &stats,
		`SELECT * FROM "stats" WHERE slot_id=$1 AND segment_id=$2`, slotID, segmentID)
	if err != nil {
		return nil, fmt.Errorf("couldn't get stats for slot_id=%d/segment_id=%d: %w", slotID, segmentID, err)
	}

	return stats, nil
}

func (sr *StatsRepository) TrackImpression(ctx context.Context, crID, slotID, segID models.ID) error {
	return sr.updateStats(ctx, impressions, crID, slotID, segID)
}

func (sr *StatsRepository) TrackConversion(ctx context.Context, crID, slotID, segID models.ID) error {
	return sr.updateStats(ctx, conversions, crID, slotID, segID)
}

func (sr *StatsRepository) updateStats(ctx context.Context, kind statsKind,
	creativeID, slotID, segmentID models.ID) error {
	query := `INSERT INTO "stats" ("$1", creative_id, slot_id, segment_id) 
				SELECT * FROM  (VALUES (1, $2, $3, $4)) AS i($1, creative_id, slot_id, segment_id)
				WHERE EXISTS (
				    SELECT FROM "slot_creative" sc
   					WHERE  sc.slot_id = i.slot_id
   					AND    sc.creative_id =i.creative_id
				)
			ON CONFLICT (creative_id, slot_id, segment_id) DO UPDATE SET "$1" = "$1" + 1`
	_, err := sr.db.ExecContext(ctx, query, kind.String(), creativeID, slotID, segmentID)
	if err != nil {
		return fmt.Errorf("couldn't update stats for creative_id=%d, slot_id=%d, segment_id=%d: %w",
			creativeID, slotID, segmentID, err)
	}

	return nil
}
