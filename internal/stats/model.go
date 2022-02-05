package stats

import (
	"time"

	"github.com/pls87/creative-rotation/internal/storage/models"
)

type Event struct {
	CreativeID models.ID `json:"creative_id"`
	SegmentID  models.ID `json:"segment_idid"`
	SlotID     models.ID `json:"slot_id"`
	Time       time.Time `json:"time"`
}
