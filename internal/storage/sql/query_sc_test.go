package sql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SCQuerySuite struct {
	QuerySuite
}

func (s *SCQuerySuite) TestCreativeSlots() {
	_, e := s.storage.Creatives().Slots(context.Background(), 0)
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`SELECT s.* FROM "slot_creative" sc
		INNER JOIN "slot" s ON sc.slot_id=s."ID" 
		WHERE sc.creative_id = $1`)
}

func (s *SCQuerySuite) TestSlotCreatives() {
	_, e := s.storage.Slots().Creatives(context.Background(), 0)
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`SELECT s.* FROM "slot_creative" sc
		INNER JOIN "creative" s ON sc.creative_id=s."ID" 
		WHERE sc.slot_id = $1`)
}

func (s *SCQuerySuite) TestAllSlotCreatives() {
	_, e := s.storage.Creatives().AllSlotCreatives(context.Background())
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`SELECT sc.slot_id, sc.creative_id, s.description as slot_desc, 
			cr.description as creative_desc 
		FROM "slot_creative" sc 
		INNER JOIN "slot" s ON sc.slot_id=s."ID"
		INNER JOIN "creative" cr ON sc.creative_id=cr."ID"`)
}

func (s *SCQuerySuite) TestAddToSlot() {
	e := s.storage.Creatives().ToSlot(context.Background(), 0, 0)
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`INSERT INTO "slot_creative" (creative_id, slot_id) VALUES ($1, $2)`)
}

func (s *SCQuerySuite) TestDeleteFromSlot() {
	e := s.storage.Creatives().FromSlot(context.Background(), 0, 0)
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`DELETE FROM "slot_creative" WHERE creative_id = $1 AND slot_id=$2`)
}

func (s *SCQuerySuite) TestExistsInSlot() {
	_, e := s.storage.Creatives().InSlot(context.Background(), 0, 0)
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`SELECT * FROM "slot_creative" WHERE creative_id = $1 AND slot_id = $2`)
}

func TestSCQueries(t *testing.T) {
	suite.Run(t, new(SCQuerySuite))
}
