package business

import (
	"testing"

	"github.com/pls87/creative-rotation/internal/storage/models"
	"github.com/stretchr/testify/require"
)

type step struct {
	action      string
	times       int
	creative    models.ID
	from        []models.ID
	zero, total int
}

var stats = []models.Stats{
	{CreativeID: 1, SlotID: 1, SegmentID: 1, Impressions: 0, Conversions: 0},
	{CreativeID: 3, SlotID: 1, SegmentID: 1, Impressions: 0, Conversions: 0},
	{CreativeID: 5, SlotID: 1, SegmentID: 1, Impressions: 0, Conversions: 0},
	{CreativeID: 9, SlotID: 1, SegmentID: 1, Impressions: 0, Conversions: 0},
	{CreativeID: 2, SlotID: 1, SegmentID: 1, Impressions: 0, Conversions: 0},
}

var steps = []step{
	{action: "imp", creative: 5, zero: 4, total: 1},
	{action: "imp", creative: 1, zero: 3, total: 2},
	{action: "imp", creative: 2, zero: 2, total: 3},
	{action: "conv", creative: 5, zero: 2, total: 3},
	{action: "next", from: []models.ID{3, 9}, zero: 2, total: 3, times: 3},
	{action: "imp", creative: 9, zero: 1, total: 4},
	{action: "next", from: []models.ID{3}, zero: 1, total: 4, times: 2},
	{action: "conv", creative: 9, zero: 1, total: 4},
	{action: "imp", creative: 3, zero: 0, total: 5},
	{action: "next", from: []models.ID{5, 9}, zero: 0, total: 5, times: 3},
	{action: "imp", creative: 3, zero: 0, total: 6},
	{action: "next", from: []models.ID{5, 9}, zero: 0, total: 6, times: 3},
	{action: "imp", creative: 1, zero: 0, total: 7},
	{action: "next", from: []models.ID{5, 9}, zero: 0, total: 7, times: 3},
	{action: "conv", creative: 1, zero: 0, total: 7},
	{action: "next", from: []models.ID{1, 5, 9}, zero: 0, total: 7, times: 5},
	{action: "conv", creative: 5, zero: 0, total: 7},
	{action: "next", from: []models.ID{5}, zero: 0, total: 7},
	{action: "imp", creative: 5, zero: 0, total: 8},
	{action: "imp", creative: 5, zero: 0, total: 9},
	{action: "next", from: []models.ID{1, 9}, zero: 0, total: 9, times: 5},
}

func imp(creative models.ID, stats []models.Stats) []models.Stats {
	for i, v := range stats {
		if v.CreativeID == creative {
			v.Impressions++
			stats[i] = v
		}
	}
	return stats
}

func conv(creative models.ID, stats []models.Stats) []models.Stats {
	for i, v := range stats {
		if v.CreativeID == creative {
			v.Conversions++
			stats[i] = v
		}
	}
	return stats
}

func find(ids []models.ID, id models.ID) bool {
	for _, i := range ids {
		if i == id {
			return true
		}
	}
	return false
}

func TestUSB1Playback(t *testing.T) {
	stats := stats
	t.Helper()
	for _, s := range steps {
		times := 1
		if s.times > 0 {
			times = s.times
		}
		for i := 0; i < times; i++ {
			switch s.action {
			case "imp":
				stats = imp(s.creative, stats)
			case "conv":
				stats = conv(s.creative, stats)
			case "next":
				c := NextCreative(stats)
				require.Truef(t, find(s.from, c), "%d wasn't found in %v", c, s.from)
			}
			zero, total := aggregate(stats)
			require.Equal(t, s.zero, len(zero))
			require.Equal(t, uint64(s.total), total)
		}
	}
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

	require.Equalf(t, expectedTotal, total, "Impressions are not aggregated correctly")
	require.Equalf(t, len(expectedZeroStats), len(zeroStats), "Expected %d empty stats but got: %d",
		len(expectedZeroStats), len(zeroStats))
	require.Equalf(t, expectedZeroStats, zeroStats, "Empty stats aggregated incorrectly")
}
