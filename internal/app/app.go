package app

import (
	"context"

	"github.com/pls87/creative-rotation/internal/business"

	"github.com/pls87/creative-rotation/internal/config"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type Application interface {
	AllCreatives(ctx context.Context) ([]models.Creative, error)
	NewCreative(ctx context.Context, c models.Creative) (created models.Creative, err error)
	AllSlots(ctx context.Context) ([]models.Slot, error)
	NewSlot(ctx context.Context, s models.Slot) (created models.Slot, err error)
	AllSegments(ctx context.Context) ([]models.Segment, error)
	NewSegment(ctx context.Context, s models.Segment) (created models.Segment, err error)
	AddToSlot(ctx context.Context, creativeId, slotId models.ID) error
	RemoveFromSlot(ctx context.Context, creativeId, slotId models.ID) error
	TrackConversion(ctx context.Context, creativeId, slotId, segmentId models.ID) error
	NextCreative(ctx context.Context, slotId, segmentId models.ID) (models.ID, error)
}

type App struct {
	logger  *logger.Logger
	storage basic.Storage
	cfg     config.Config
}

func (a *App) AllCreatives(ctx context.Context) (created []models.Creative, err error) {
	return a.storage.Creatives().All(ctx)
}

func (a *App) NewCreative(ctx context.Context, c models.Creative) (created models.Creative, err error) {
	return a.storage.Creatives().Create(ctx, c)
}

func (a *App) AllSlots(ctx context.Context) (created []models.Slot, err error) {
	return a.storage.Slots().All(ctx)
}

func (a *App) NewSlot(ctx context.Context, s models.Slot) (created models.Slot, err error) {
	return a.storage.Slots().Create(ctx, s)
}

func (a *App) AllSegments(ctx context.Context) (created []models.Segment, err error) {
	return a.storage.Segments().All(ctx)
}

func (a *App) NewSegment(ctx context.Context, s models.Segment) (created models.Segment, err error) {
	return a.storage.Segments().Create(ctx, s)
}

func (a *App) AddToSlot(ctx context.Context, creativeId, slotId models.ID) error {
	return a.storage.Creatives().ToSlot(ctx, creativeId, slotId)
}

func (a *App) RemoveFromSlot(ctx context.Context, creativeId, slotId models.ID) error {
	return a.storage.Creatives().FromSlot(ctx, creativeId, slotId)
}

func (a *App) TrackConversion(ctx context.Context, creativeId, slotId, segmentId models.ID) error {
	return a.storage.Creatives().TrackConversion(ctx, creativeId, slotId, segmentId)
}

func (a *App) NextCreative(ctx context.Context, slotId, segmentId models.ID) (models.ID, error) {
	stats, err := a.storage.Stats().StatsSlotSegment(ctx, slotId, segmentId)
	if err != nil {
		a.logger.WithContext(ctx).Errorf("Next creative: %s", err)
		return 0, err
	}
	return business.NextCreative(stats), nil
}

func New(logger *logger.Logger, storage basic.Storage, cfg config.Config) Application {
	return &App{logger, storage, cfg}
}
