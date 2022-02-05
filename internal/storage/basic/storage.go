package basic

import (
	"context"
	"errors"

	"github.com/pls87/creative-rotation/internal/storage/models"
)

var (
	ErrDoesNotExist          = errors.New("entity doesn't exist")
	ErrCreativeNotInSlot     = errors.New("creative isn't attached to slot")
	ErrCreativeAlreadyInSlot = errors.New("creative already in slot")
)

type Storage interface {
	Creatives() CreativeRepository
	Slots() SlotRepository
	Segments() SegmentRepository
	Stats() StatsRepository
	Init(ctx context.Context) error
	Dispose() error
}

type SegmentRepository interface {
	All(ctx context.Context) ([]models.Segment, error)
	Create(ctx context.Context, s models.Segment) (added models.Segment, err error)
	Delete(ctx context.Context, id models.ID) error
}

type SlotRepository interface {
	All(ctx context.Context) ([]models.Slot, error)
	Create(ctx context.Context, s models.Slot) (added models.Slot, err error)
	Delete(ctx context.Context, id models.ID) error
}

type CreativeRepository interface {
	All(ctx context.Context) ([]models.Creative, error)
	Create(ctx context.Context, c models.Creative) (added models.Creative, err error)
	Delete(ctx context.Context, id models.ID) error
	ToSlot(ctx context.Context, creativeID, slotID models.ID) error
	FromSlot(ctx context.Context, creativeID, slotID models.ID) error
	InSlot(ctx context.Context, creativeID, slotID models.ID) (bool, error)
	TrackImpression(ctx context.Context, impression models.Impression) error
	TrackConversion(ctx context.Context, conversion models.Conversion) error
}

type StatsRepository interface {
	StatsSlotSegment(ctx context.Context, slotID, segmentID models.ID) ([]models.Stats, error)
	UpdateStatsImpression(ctx context.Context, impression models.Impression) error
	UpdateStatsConversion(ctx context.Context, conversion models.Conversion) error
}
