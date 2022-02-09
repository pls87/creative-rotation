package sql

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
)

type StatQuerySuite struct {
	QuerySuite
}

func (s *StatQuerySuite) TestStatForSlotSegment() {
	_, e := s.storage.Stats().StatsSlotSegment(context.Background(), 0, 0)
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`SELECT * FROM "stats" WHERE slot_id=$1 AND segment_id=$2`)
}

func (s *StatQuerySuite) TestTrackImpression() {
	e := s.storage.Stats().TrackImpression(context.Background(), 0, 0, 0)
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`INSERT INTO "stats" (impressions, creative_id, slot_id, segment_id) 
		SELECT * FROM  (VALUES (1, $1, $2, $3)) AS t(impressions, creative_id, slot_id, segment_id)
		WHERE EXISTS (
			SELECT FROM "slot_creative" sc
			WHERE  sc.slot_id = t.slot_id
				AND    sc.creative_id = t.creative_id
		)
		ON CONFLICT (creative_id, slot_id, segment_id) DO UPDATE SET impressions = EXCLUDED.impressions + 1`)
}

func (s *StatQuerySuite) TestTrackConversion() {
	e := s.storage.Stats().TrackConversion(context.Background(), 0, 0, 0)
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`INSERT INTO "stats" (conversions, creative_id, slot_id, segment_id) 
		SELECT * FROM  (VALUES (1, $1, $2, $3)) AS t(conversions, creative_id, slot_id, segment_id)
		WHERE EXISTS (
			SELECT FROM "slot_creative" sc
			WHERE  sc.slot_id = t.slot_id
				AND    sc.creative_id = t.creative_id
		)
		ON CONFLICT (creative_id, slot_id, segment_id) DO UPDATE SET conversions = EXCLUDED.conversions + 1`)
}

func TestStatQueries(t *testing.T) {
	suite.Run(t, new(StatQuerySuite))
}
