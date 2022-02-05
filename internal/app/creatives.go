package app

import (
	"context"

	"github.com/pls87/creative-rotation/internal/business"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/stats"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type CreativeApplication interface {
	All(ctx context.Context) ([]models.Creative, error)
	New(ctx context.Context, c models.Creative) (created models.Creative, err error)
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

func (a *CreativeApp) AddToSlot(ctx context.Context, creativeID, slotID models.ID) error {
	return a.storage.Creatives().ToSlot(ctx, creativeID, slotID)
}

func (a *CreativeApp) RemoveFromSlot(ctx context.Context, creativeID, slotID models.ID) error {
	return a.storage.Creatives().FromSlot(ctx, creativeID, slotID)
}

func (a *CreativeApp) TrackConversion(ctx context.Context, conversion models.Conversion) error {
	return a.stats.Produce("conversion", stats.Event{
		CreativeID: conversion.CreativeID,
		SegmentID:  conversion.SegmentID,
		SlotID:     conversion.SlotID,
		Time:       conversion.Time,
	})
}

func (a *CreativeApp) TrackImpression(ctx context.Context, impression models.Impression) error {
	return a.stats.Produce("conversion", stats.Event{
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
		a.logger.WithContext(ctx).Errorf("couldn't get stats to calculate next creative: %s", err)
	} else {
		next.ID = business.NextCreative(stat)
	}

	return next, err
}
