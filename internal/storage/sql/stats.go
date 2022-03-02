package sql

import (
	"context"

	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
	"github.com/pls87/creative-rotation/internal/storage/sql/errors"
	"github.com/pls87/creative-rotation/internal/storage/sql/queries"
)

type StatsRepository struct {
	db DB
}

func (sr *StatsRepository) StatsSlotSegment(ctx context.Context, slotID, segmentID models.ID) ([]models.Stats, error) {
	var stats []models.Stats
	err := sr.db.SelectContext(ctx, &stats, queries.Stats.For(), slotID, segmentID)
	if err != nil {
		return nil, errors.Stats.For(slotID, segmentID, err)
	}

	return stats, nil
}

func (sr *StatsRepository) TrackImpression(ctx context.Context, crID, slotID, segID models.ID) error {
	return sr.updateStats(ctx, queries.ImpressionField, crID, slotID, segID)
}

func (sr *StatsRepository) TrackConversion(ctx context.Context, crID, slotID, segID models.ID) error {
	return sr.updateStats(ctx, queries.ConversionField, crID, slotID, segID)
}

func (sr *StatsRepository) updateStats(ctx context.Context, field string, crID, slotID, segID models.ID) error {
	res, err := sr.db.ExecContext(ctx, queries.Stats.Track(field, crID, slotID, segID))
	if err != nil {
		return errors.Stats.Track(field, crID, slotID, segID, err)
	}

	if affected, _ := res.RowsAffected(); affected == 0 {
		return errors.Stats.Track(field, crID, slotID, segID, basic.ErrCreativeNotInSlot)
	}

	return nil
}
