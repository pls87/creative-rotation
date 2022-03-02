package helpers

import "github.com/pls87/creative-rotation/internal/storage/models"

type CreativeCollection struct {
	Creatives []models.Creative `json:"creatives"`
}

type SlotCollection struct {
	Slots []models.Slot `json:"slots"`
}

type SegmentCollection struct {
	Segments []models.Segment `json:"segments"`
}

type SCCollection struct {
	SC []models.SlotCreative `json:"slot_creatives"`
}
