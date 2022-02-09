//go:build integration
// +build integration

package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/pls87/creative-rotation/internal/tests/integration/helpers"
	"github.com/stretchr/testify/suite"
)

type StatsSuite struct {
	BaseSuite
	entitiesH  *helpers.EntityHelper
	creativesH *helpers.CreativeHelper
	statsH     *helpers.StatHelper
}

func (s *StatsSuite) SetupSuite() {
	s.BaseSuite.SetupSuite()
	s.entitiesH = helpers.NewEntityHelper(s.baseURL)
	s.creativesH = helpers.NewCreativeHelper(s.baseURL)
	s.statsH = helpers.NewStatHelper(s.baseURL)
}

func (s *StatsSuite) TestStatsUpdate() {
	cr1 := s.entitiesH.New(s.T(), "creative", gofakeit.BuzzWord())
	cr2 := s.entitiesH.New(s.T(), "creative", gofakeit.BuzzWord())
	sl1 := s.entitiesH.New(s.T(), "slot", gofakeit.DomainName())
	sl2 := s.entitiesH.New(s.T(), "slot", gofakeit.DomainName())
	seg1 := s.entitiesH.New(s.T(), "segment", gofakeit.JobTitle())
	seg2 := s.entitiesH.New(s.T(), "segment", gofakeit.JobTitle())

	s.addToSlot(cr2.ID, sl1.ID)
	s.trackImpression(cr2.ID, sl1.ID, seg1.ID)
	time.Sleep(100 * time.Millisecond)
	stats := s.getStats(cr2.ID, sl1.ID, seg1.ID)

	s.Equalf(helpers.Stats{
		SlotID: sl1.ID, CreativeID: cr2.ID, SegmentID: seg1.ID,
		Impressions: 1,
		Conversions: 0,
	}, stats, "Impression number wasn't updated")

	s.addToSlot(cr1.ID, sl2.ID)
	s.trackConversion(cr1.ID, sl2.ID, seg2.ID)

	time.Sleep(100 * time.Millisecond)
	stats = s.getStats(cr1.ID, sl2.ID, seg2.ID)
	s.Equalf(helpers.Stats{
		SlotID: sl2.ID, CreativeID: cr1.ID, SegmentID: seg2.ID,
		Impressions: 0,
		Conversions: 1,
	}, stats, "Conversion number wasn't updated")
}

func (s *StatsSuite) addToSlot(creativeID, slotID int) {
	code, _, err := s.creativesH.AddToSlot(creativeID, slotID)
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equalf(http.StatusOK, code, "status %d expected but was %d", http.StatusOK, code)
}

func (s *StatsSuite) trackImpression(creativeID, slotID, segmentID int) {
	code, _, err := s.statsH.TrackImpression(creativeID, slotID, segmentID)
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equalf(http.StatusOK, code, "status %d expected but was %d", http.StatusOK, code)
}

func (s *StatsSuite) trackConversion(creativeID, slotID, segmentID int) {
	code, _, err := s.statsH.TrackConversion(creativeID, slotID, segmentID)
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equalf(http.StatusOK, code, "status %d expected but was %d", http.StatusOK, code)
}

func (s *StatsSuite) getStats(creativeID, slotID, segmentID int) (stats helpers.Stats) {
	err := s.client.Db.Get(&stats, `SELECT * FROM "stats" WHERE creative_id=$1 AND slot_id=$2 AND segment_id=$3`,
		creativeID, slotID, segmentID)
	s.NoError(err)
	return stats
}

func TestStats(t *testing.T) {
	suite.Run(t, new(StatsSuite))
}
