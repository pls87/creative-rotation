package sql

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"strings"
	"sync"
	"testing"

	"github.com/pls87/creative-rotation/internal/config"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/pls87/creative-rotation/internal/storage/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var ErrEmpty = errors.New("")

type TestDB struct {
	sync.Mutex
	*testing.T
	lastQuery string
}

func (db *TestDB) CheckLastQuery(expected string) {
	defer db.Unlock()
	expectedQ := strings.ReplaceAll(strings.ToLower(expected), "\n", " ")
	actualQ := strings.ReplaceAll(strings.ToLower(db.lastQuery), "\n", " ")
	space := regexp.MustCompile(`\s+`)
	space.ReplaceAllString(expectedQ, " ")
	require.Equalf(db, space.ReplaceAllString(expectedQ, " "), space.ReplaceAllString(actualQ, " "),
		"Queries are not matched")
}

func (db *TestDB) SetLastQuery(query string) {
	db.Lock()
	db.lastQuery = query
}

func (db *TestDB) Open(ctx context.Context, cfg config.DBConf) error {
	return nil
}

func (db *TestDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	db.SetLastQuery(query)
	return nil, ErrEmpty
}

func (db *TestDB) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	db.SetLastQuery(query)
	return ErrEmpty
}

func (db *TestDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	db.SetLastQuery(query)
	return nil, ErrEmpty
}

func (db *TestDB) InsertRow(ctx context.Context, query string, args ...interface{}) (models.ID, error) {
	db.SetLastQuery(query)

	return 0, ErrEmpty
}

func (db *TestDB) Close() error {
	return nil
}

type QuerySuite struct {
	suite.Suite
	storage basic.Storage
	testDB  *TestDB
}

func (s *QuerySuite) SetupSuite() {
	s.testDB = &TestDB{T: s.T()}
	s.storage = New(config.DBConf{}, s.testDB)
	s.NoError(s.storage.Init(context.Background()))
}
