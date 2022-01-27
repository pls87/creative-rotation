package app

import (
	"context"

	"github.com/pls87/creative-rotation/internal/business"

	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/storage/basic"

	"github.com/pls87/creative-rotation/internal/storage/models"
)

type CreativeApplication interface {
	All(ctx context.Context) ([]models.Creative, error)
	New(ctx context.Context, c models.Creative) (created models.Creative, err error)
	AddToSlot(ctx context.Context, creativeId, slotId models.ID) error
	RemoveFromSlot(ctx context.Context, creativeId, slotId models.ID) error
	TrackConversion(ctx context.Context, conversion models.Conversion) error
	Next(ctx context.Context, slotId, segmentId models.ID) (models.Creative, error)
}

type CreativeApp struct {
	logger  *logger.Logger
	storage basic.Storage
}

func (a *CreativeApp) All(ctx context.Context) (collection []models.Creative, err error) {
	return a.storage.Creatives().All(ctx)
}

func (a *CreativeApp) New(ctx context.Context, c models.Creative) (created models.Creative, err error) {
	return a.storage.Creatives().Create(ctx, c)
}

func (a *CreativeApp) AddToSlot(ctx context.Context, creativeId, slotId models.ID) error {
	return a.storage.Creatives().ToSlot(ctx, creativeId, slotId)
}

func (a *CreativeApp) RemoveFromSlot(ctx context.Context, creativeId, slotId models.ID) error {
	return a.storage.Creatives().FromSlot(ctx, creativeId, slotId)
}

func (a *CreativeApp) TrackConversion(ctx context.Context, conversion models.Conversion) error {
	return a.storage.Creatives().TrackConversion(ctx, conversion)
}

func (a *CreativeApp) Next(ctx context.Context, slotId, segmentId models.ID) (models.Creative, error) {
	next := models.Creative{}
	stats, err := a.storage.Stats().StatsSlotSegment(ctx, slotId, segmentId)
	if err != nil {
		a.logger.WithContext(ctx).Errorf("Next creative: %s", err)
	} else {
		next.ID = business.NextCreative(stats)
	}

	return next, err
}
