package business

import (
	"math"
	"math/rand"

	"github.com/pls87/creative-rotation/internal/storage/models"
)

func NextCreative(stats []models.Stats) (creativeID models.ID) {
	zeroImp, totalImp := aggregate(stats)
	if len(zeroImp) > 0 {
		return zeroImp[rand.Intn(len(zeroImp))].CreativeID //nolint: gosec
	}
	var cur float64
	max := math.Inf(-1)
	for _, s := range stats {
		cur = valueToMaximize(s, totalImp)
		if cur > max {
			max = cur
			creativeID = s.CreativeID
		}
	}
	return creativeID
}

func valueToMaximize(stats models.Stats, totalImp uint64) float64 {
	impressions := float64(stats.Impressions)
	totalImpressions := float64(totalImp)
	conversions := float64(stats.Conversions)

	return conversions/impressions + math.Sqrt(2*math.Log(totalImpressions)/impressions)
}

func aggregate(stats []models.Stats) (zeroStats []models.Stats, totalImp uint64) {
	zeroStats = make([]models.Stats, 0, 10)
	for _, s := range stats {
		totalImp += s.Impressions
		if s.Impressions == 0 {
			zeroStats = append(zeroStats, s)
		}
	}

	return zeroStats, totalImp
}
