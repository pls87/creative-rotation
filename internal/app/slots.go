package app

import (
	"context"

	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type SlotApplication interface {
	All(ctx context.Context) ([]models.Slot, error)
	New(ctx context.Context, s models.Slot) (created models.Slot, err error)
}

type SlotApp struct {
	logger  *logger.Logger
	storage basic.Storage
}

func (a *SlotApp) All(ctx context.Context) (created []models.Slot, err error) {
	return a.storage.Slots().All(ctx)
}

func (a *SlotApp) New(ctx context.Context, s models.Slot) (created models.Slot, err error) {
	return a.storage.Slots().Create(ctx, s)
}
