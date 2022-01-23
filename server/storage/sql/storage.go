package sql

import (
	"context"

	"github.com/jmoiron/sqlx"
	// init postgres driver.
	_ "github.com/lib/pq"
	"github.com/pls87/creative-rotation/server/config"
	"github.com/pls87/creative-rotation/server/storage/basic"
)

type Storage struct {
	cfg      config.DBConf
	db       *sqlx.DB
	segments *SegmentRepository
	slots    *SlotRepository
}

func New(cfg config.DBConf) *Storage {
	return &Storage{
		segments: &SegmentRepository{},
		slots:    &SlotRepository{},
		cfg:      cfg,
	}
}

func (s *Storage) Segments() basic.SegmentRepository {
	return s.segments
}

func (s *Storage) Slots() basic.SlotRepository {
	return s.slots
}

func (s *Storage) Init(ctx context.Context) error {
	db, err := sqlx.ConnectContext(ctx, "postgres", s.cfg.ConnString())
	if err == nil {
		s.db = db
		s.segments.Init(ctx)
		s.segments.db = s.db
	}
	return err
}

func (s *Storage) Dispose() error {
	return s.db.Close()
}
