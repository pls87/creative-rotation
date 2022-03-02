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
	s.testDB.CheckLastQuery(`SELECT sc.slot_id as slot_id, 
		sc.creative_id as creative_id, seg."ID" as segment_id,
		CASE WHEN st.impressions is NULL THEN 0 ELSE st.impressions END as impressions,
		CASE WHEN st.conversions is NULL THEN 0 ELSE st.conversions END as conversions
		FROM "slot_creative" sc CROSS JOIN "segment" seg LEFT JOIN "stats" st ON sc.creative_id = st.creative_id 
		AND sc.slot_id = st.slot_id AND seg."ID" = st.segment_id
		WHERE sc.slot_id = $1 and seg."ID" = $2`)
}

func (s *StatQuerySuite) TestTrackImpression() {
	e := s.storage.Stats().TrackImpression(context.Background(), 0, 0, 0)
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`INSERT INTO "stats" (impressions, creative_id, slot_id, segment_id) 
		VALUES (1, 0, 0, 0)
	ON CONFLICT (creative_id, slot_id, segment_id) 
	DO UPDATE SET impressions = (SELECT impressions FROM "stats" WHERE creative_id=0 AND slot_id=0 AND segment_id=0) + 1`)
}

func (s *StatQuerySuite) TestTrackConversion() {
	e := s.storage.Stats().TrackConversion(context.Background(), 0, 0, 0)
	s.ErrorIs(e, ErrEmpty)
	s.testDB.CheckLastQuery(`INSERT INTO "stats" (conversions, creative_id, slot_id, segment_id) 
		VALUES (1, 0, 0, 0)
	ON CONFLICT (creative_id, slot_id, segment_id) 
	DO UPDATE SET conversions = (SELECT conversions FROM "stats" WHERE creative_id=0 AND slot_id=0 AND segment_id=0) + 1`)
}

func TestStatQueries(t *testing.T) {
	suite.Run(t, new(StatQuerySuite))
}
