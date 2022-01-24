package app

import (
	"github.com/pls87/creative-rotation/internal/config"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/storage/basic"
)

type Application interface{}

type App struct {
	logger  *logger.Logger
	storage basic.Storage
	cfg     config.Config
}

func New(logger *logger.Logger, storage basic.Storage, cfg config.Config) Application {
	return &App{logger, storage, cfg}
}
