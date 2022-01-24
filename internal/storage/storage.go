package storage

import (
	"github.com/pls87/creative-rotation/internal/config"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/sql"
)

func New(cfg config.DBConf) basic.Storage {
	return sql.New(cfg)
}
