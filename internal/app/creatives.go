package app

import (
	"context"
	"fmt"
	"time"

	"github.com/pls87/creative-rotation/internal/business"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/stats"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type CreativeApplication interface {
	All(ctx context.Context) ([]models.Creative, error)
	New(ctx context.Context, c models.Creative) (created models.Creative, err error)
	Slots(ctx context.Context, creativeID models.ID) ([]models.Slot, error)
	AllSlotCreatives(ctx context.Context) ([]models.SlotCreative, error)
	AddToSlot(ctx context.Context, creativeID, slotID models.ID) error
	RemoveFromSlot(ctx context.Context, creativeID, slotID models.ID) error
	TrackConversion(ctx context.Context, conversion models.Conversion) error
	TrackImpression(ctx context.Context, conversion models.Impression) error
	Next(ctx context.Context, slotID, segmentID models.ID) (models.Creative, error)
}

type CreativeApp struct {
	logger  *logger.Logger
	storage basic.Storage
	stats   stats.Producer
}

func (a *CreativeApp) All(ctx context.Context) (collection []models.Creative, err error) {
	return a.storage.Creatives().All(ctx)
}

func (a *CreativeApp) New(ctx context.Context, c models.Creative) (created models.Creative, err error) {
	return a.storage.Creatives().Create(ctx, c)
}

func (a *CreativeApp) Slots(ctx context.Context, creativeID models.ID) ([]models.Slot, error) {
	return a.storage.Creatives().Slots(ctx, creativeID)
}

func (a *CreativeApp) AllSlotCreatives(ctx context.Context) ([]models.SlotCreative, error) {
	return a.storage.Creatives().AllSlotCreatives(ctx)
}

func (a *CreativeApp) AddToSlot(ctx context.Context, creativeID, slotID models.ID) error {
	return a.storage.Creatives().ToSlot(ctx, creativeID, slotID)
}

func (a *CreativeApp) RemoveFromSlot(ctx context.Context, creativeID, slotID models.ID) error {
	return a.storage.Creatives().FromSlot(ctx, creativeID, slotID)
}

func (a *CreativeApp) TrackConversion(_ context.Context, conversion models.Conversion) error {
	if conversion.Time.IsZero() {
		conversion.Time = time.Now()
	}
	return a.stats.Produce(stats.ConversionKey, stats.Event{
		CreativeID: conversion.CreativeID,
		SegmentID:  conversion.SegmentID,
		SlotID:     conversion.SlotID,
		Time:       conversion.Time,
	})
}

func (a *CreativeApp) TrackImpression(_ context.Context, impression models.Impression) error {
	if impression.Time.IsZero() {
		impression.Time = time.Now()
	}
	return a.stats.Produce(stats.ImpressionKey, stats.Event{
		CreativeID: impression.CreativeID,
		SegmentID:  impression.SegmentID,
		SlotID:     impression.SlotID,
		Time:       impression.Time,
	})
}

func (a *CreativeApp) Next(ctx context.Context, slotID, segmentID models.ID) (models.Creative, error) {
	next := models.Creative{}
	stat, err := a.storage.Stats().StatsSlotSegment(ctx, slotID, segmentID)
	if err != nil {
		return next, fmt.Errorf("couldn't get stats to calculate next creative: %w", err)
	}
	id, e := business.NextCreative(stat)
	if e != nil {
		return next, fmt.Errorf("couldn't calculate next creative: %w", e)
	}
	next.ID = id

	return next, nil
}
