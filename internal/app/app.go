package app

import (
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/stats"
	"github.com/pls87/creative-rotation/internal/storage/basic"
)

type Application interface {
	Creatives() CreativeApplication
	Slots() SlotApplication
	Segments() SegmentApplication
}

type App struct {
	creatives CreativeApp
	slots     SlotApp
	segments  SegmentApp
}

func (a *App) Creatives() CreativeApplication {
	return &a.creatives
}

func (a *App) Slots() SlotApplication {
	return &a.slots
}

func (a *App) Segments() SegmentApplication {
	return &a.segments
}

func New(logger *logger.Logger, storage basic.Storage, stats stats.Producer) Application {
	return &App{
		creatives: CreativeApp{
			logger:  logger,
			storage: storage,
			stats:   stats,
		},
		slots: SlotApp{
			logger:  logger,
			storage: storage,
		},
		segments: SegmentApp{
			logger:  logger,
			storage: storage,
		},
	}
}
