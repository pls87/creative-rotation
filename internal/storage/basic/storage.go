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
	Init(ctx context.Context) error
	All(ctx context.Context) ([]models.Segment, error)
	Create(ctx context.Context, s models.Segment) (added models.Segment, err error)
	Delete(ctx context.Context, id models.ID) error
}

type SlotRepository interface {
	Init(ctx context.Context) error
	All(ctx context.Context) ([]models.Slot, error)
	Create(ctx context.Context, s models.Slot) (added models.Slot, err error)
	Delete(ctx context.Context, id models.ID) error
}

type CreativeRepository interface {
	Init(ctx context.Context) error
	All(ctx context.Context) ([]models.Creative, error)
	Create(ctx context.Context, c models.Creative) (added models.Creative, err error)
	Delete(ctx context.Context, id models.ID) error
	ToSlot(ctx context.Context, creativeId, slotId models.ID) error
	FromSlot(ctx context.Context, creativeId, slotId models.ID) error
	InSlot(ctx context.Context, creativeId, slotId models.ID) (bool, error)
	TrackImpression(ctx context.Context, impression models.Impression) error
	TrackConversion(ctx context.Context, conversion models.Conversion) error
}

type StatsRepository interface {
	AllStats(ctx context.Context) ([]models.Stats, error)
	StatsSlotSegment(ctx context.Context, slotId, segmentId models.ID) ([]models.Stats, error)
}
