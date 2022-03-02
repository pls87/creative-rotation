package errors

import (
	"fmt"

	"github.com/pls87/creative-rotation/internal/storage/models"
)

type StatsErrors struct{}

func (c *StatsErrors) For(slotID, segmentID models.ID, err error) error {
	return fmt.Errorf("couldn't get stats for slot_id=%d/segment_id=%d: %w", slotID, segmentID, err)
}

func (c *StatsErrors) Track(field string, creativeID, slotID, segmentID models.ID, err error) error {
	return fmt.Errorf("couldn't track %s for creative_id=%d, slot_id=%d, segment_id=%d: %w",
		field, creativeID, slotID, segmentID, err)
}

var Stats = StatsErrors{}
