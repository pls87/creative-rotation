package business

import (
	"testing"

	"github.com/pls87/creative-rotation/internal/storage/models"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type step struct {
	action      string
	times       int
	creative    models.ID
	from        []models.ID
	zero, total int
}

var steps = []step{
	{action: "imp", creative: 5, zero: 4, total: 1},
	{action: "imp", creative: 1, zero: 3, total: 2},
	{action: "imp", creative: 2, zero: 2, total: 3},
	{action: "conv", creative: 5, zero: 2, total: 3},
	// choose creative of 3 and 9 because they have zero stats
	{action: "next", from: []models.ID{3, 9}, zero: 2, total: 3, times: 3},
	{action: "imp", creative: 9, zero: 1, total: 4},
	// choose creative 9 because it has zero stats
	{action: "next", from: []models.ID{3}, zero: 1, total: 4, times: 2},
	{action: "conv", creative: 9, zero: 1, total: 4},
	{action: "imp", creative: 3, zero: 0, total: 5},
	// choose creative of 5 and 9 because they have same stats
	{action: "next", from: []models.ID{5, 9}, zero: 0, total: 5, times: 3},
	{action: "imp", creative: 3, zero: 0, total: 6},
	// choose creative of 5 and 9 because they have same stats
	{action: "next", from: []models.ID{5, 9}, zero: 0, total: 6, times: 3},
	{action: "imp", creative: 1, zero: 0, total: 7},
	// choose creative of 5 and 9 because they have same stats
	{action: "next", from: []models.ID{5, 9}, zero: 0, total: 7, times: 3},
	{action: "conv", creative: 1, zero: 0, total: 7},
	// choose creative of 1, 5 and 9 because they have same stats
	{action: "next", from: []models.ID{1, 5, 9}, zero: 0, total: 7, times: 5},
	{action: "conv", creative: 5, zero: 0, total: 7},
	// choose creative 5 because it has the best stats
	{action: "next", from: []models.ID{5}, zero: 0, total: 7},
	{action: "imp", creative: 5, zero: 0, total: 8},
	{action: "imp", creative: 5, zero: 0, total: 9},
	// choose creative of 1, 9 because they have same stats after creative 5 missed 2 impressions
	{action: "next", from: []models.ID{1, 9}, zero: 0, total: 9, times: 5},
}

func find(ids []models.ID, id models.ID) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}

type BusinessSuite struct {
	suite.Suite
	stats []models.Stats
}

func (bs *BusinessSuite) SetupTest() {
	bs.stats = []models.Stats{
		{CreativeID: 1, SlotID: 1, SegmentID: 1, Impressions: 0, Conversions: 0},
		{CreativeID: 3, SlotID: 1, SegmentID: 1, Impressions: 0, Conversions: 0},
		{CreativeID: 5, SlotID: 1, SegmentID: 1, Impressions: 0, Conversions: 0},
		{CreativeID: 9, SlotID: 1, SegmentID: 1, Impressions: 0, Conversions: 0},
		{CreativeID: 2, SlotID: 1, SegmentID: 1, Impressions: 0, Conversions: 0},
	}
}

func (bs *BusinessSuite) imp(creative models.ID) {
	for i, v := range bs.stats {
		if v.CreativeID == creative {
			v.Impressions++
			bs.stats[i] = v
		}
	}
}

func (bs *BusinessSuite) conv(creative models.ID) {
	for i, v := range bs.stats {
		if v.CreativeID == creative {
			v.Conversions++
			bs.stats[i] = v
		}
	}
}

func (bs *BusinessSuite) TestUSB1Playback() {
	stats := bs.stats
	for _, s := range steps {
		times := 1
		if s.times > 0 {
			times = s.times
		}
		for i := 0; i < times; i++ {
			switch s.action {
			case "imp":
				bs.imp(s.creative)
			case "conv":
				bs.conv(s.creative)
			case "next":
				c, e := NextCreative(stats)
				bs.NoError(e)
				bs.Truef(find(s.from, c), "%d wasn't found in %v", c, s.from)
			}
			zero, total := aggregate(stats)
			bs.Equal(s.zero, len(zero))
			bs.Equal(uint64(s.total), total)
		}
	}
}

func TestUSB1(t *testing.T) {
	suite.Run(t, new(BusinessSuite))
}

func TestEmptyStats(t *testing.T) {
	s := make([]models.Stats, 0, 1)
	_, err := NextCreative(s)
	require.ErrorIs(t, err, ErrEmptyStats)
}

func TestAggregate(t *testing.T) {
	stats := []models.Stats{
		{CreativeID: 1, SlotID: 1, SegmentID: 1, Impressions: 10, Conversions: 3},
		{CreativeID: 3, SlotID: 1, SegmentID: 1, Impressions: 7, Conversions: 1},
		{CreativeID: 5, SlotID: 1, SegmentID: 1, Impressions: 0, Conversions: 0},
		{CreativeID: 9, SlotID: 1, SegmentID: 1, Impressions: 0, Conversions: 0},
		{CreativeID: 2, SlotID: 1, SegmentID: 1, Impressions: 5, Conversions: 2},
	}
	expectedZeroStats := []models.Stats{
		{CreativeID: 5, SlotID: 1, SegmentID: 1, Impressions: 0, Conversions: 0},
		{CreativeID: 9, SlotID: 1, SegmentID: 1, Impressions: 0, Conversions: 0},
	}

	var expectedTotal uint64 = 22

	zeroStats, total := aggregate(stats)

	require.Equalf(t, expectedTotal, total, "impressions are not aggregated correctly")
	require.Equalf(t, len(expectedZeroStats), len(zeroStats), "expected %d empty stats but got: %d",
		len(expectedZeroStats), len(zeroStats))
	require.Equalf(t, expectedZeroStats, zeroStats, "empty stats aggregated incorrectly")
}
