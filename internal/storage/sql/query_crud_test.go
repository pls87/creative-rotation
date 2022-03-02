package sql

import (
	"context"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/pls87/creative-rotation/internal/storage/models"
	"github.com/stretchr/testify/suite"
)

type CRUDQuerySuite struct {
	QuerySuite
}

func (s *CRUDQuerySuite) TestNewCreative() {
	_, e := s.storage.Creatives().Create(context.Background(), models.Creative{Desc: gofakeit.BuzzWord()})
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`INSERT INTO "creative" ("description") VALUES ($1) RETURNING "ID"`)
}

func (s *CRUDQuerySuite) TestDeleteCreative() {
	e := s.storage.Creatives().Delete(context.Background(), 0)
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`DELETE FROM "creative" WHERE "ID"=$1`)
}

func (s *CRUDQuerySuite) TestAllCreative() {
	_, e := s.storage.Creatives().All(context.Background())
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`SELECT * FROM "creative"`)
}

func (s *CRUDQuerySuite) TestNewSegment() {
	_, e := s.storage.Segments().Create(context.Background(), models.Segment{Desc: gofakeit.BuzzWord()})
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`INSERT INTO "segment" ("description") VALUES ($1) RETURNING "ID"`)
}

func (s *CRUDQuerySuite) TestDeleteSegment() {
	e := s.storage.Segments().Delete(context.Background(), 0)
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`DELETE FROM "segment" WHERE "ID"=$1`)
}

func (s *CRUDQuerySuite) TestAllSegment() {
	_, e := s.storage.Segments().All(context.Background())
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`SELECT * FROM "segment"`)
}

func (s *CRUDQuerySuite) TestNewSlot() {
	_, e := s.storage.Slots().Create(context.Background(), models.Slot{Desc: gofakeit.BuzzWord()})
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`INSERT INTO "slot" ("description") VALUES ($1) RETURNING "ID"`)
}

func (s *CRUDQuerySuite) TestDeleteSlot() {
	e := s.storage.Slots().Delete(context.Background(), 0)
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`DELETE FROM "slot" WHERE "ID"=$1`)
}

func (s *CRUDQuerySuite) TestAllSlot() {
	_, e := s.storage.Slots().All(context.Background())
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`SELECT * FROM "slot"`)
}

func TestCRUDQueries(t *testing.T) {
	suite.Run(t, new(CRUDQuerySuite))
}
