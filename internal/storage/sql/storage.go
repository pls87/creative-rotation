package sql

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	// init postgres driver.
	_ "github.com/lib/pq"
	"github.com/pls87/creative-rotation/internal/config"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type DB interface {
	Open(ctx context.Context, cfg config.DBConf) error
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	InsertRow(ctx context.Context, query string, args ...interface{}) (models.ID, error)
	Close() error
}

type Database struct {
	*sqlx.DB
}

func (db *Database) Open(ctx context.Context, cfg config.DBConf) error {
	d, err := sqlx.ConnectContext(ctx, "postgres", cfg.ConnString())
	if err == nil {
		db.DB = d
	}

	return err
}

func (db *Database) InsertRow(ctx context.Context, query string, args ...interface{}) (models.ID, error) {
	var lastInsertID int
	err := db.QueryRowContext(ctx, query, args...).Scan(&lastInsertID)

	return models.ID(lastInsertID), err
}

type Storage struct {
	cfg       config.DBConf
	db        DB
	segments  *SegmentRepository
	slots     *SlotRepository
	creatives *CreativeRepository
	stats     *StatsRepository
}

func New(cfg config.DBConf, db DB) *Storage {
	if db == nil {
		db = &Database{}
	}

	return &Storage{
		segments:  &SegmentRepository{db},
		slots:     &SlotRepository{db},
		creatives: &CreativeRepository{db},
		stats:     &StatsRepository{db},
		cfg:       cfg,
		db:        db,
	}
}

func (s *Storage) Segments() basic.SegmentRepository {
	return s.segments
}

func (s *Storage) Slots() basic.SlotRepository {
	return s.slots
}

func (s *Storage) Stats() basic.StatsRepository {
	return s.stats
}

func (s *Storage) Creatives() basic.CreativeRepository {
	return s.creatives
}

func (s *Storage) Init(ctx context.Context) error {
	return s.db.Open(ctx, s.cfg)
}

func (s *Storage) Dispose() error {
	return s.db.Close()
}
