package app

import (
	"context"

	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type SegmentApplication interface {
	All(ctx context.Context) ([]models.Segment, error)
	New(ctx context.Context, s models.Segment) (created models.Segment, err error)
}

type SegmentApp struct {
	logger  *logger.Logger
	storage basic.Storage
}

func (a *SegmentApp) All(ctx context.Context) (created []models.Segment, err error) {
	return a.storage.Segments().All(ctx)
}

func (a *SegmentApp) New(ctx context.Context, s models.Segment) (created models.Segment, err error) {
	return a.storage.Segments().Create(ctx, s)
}
