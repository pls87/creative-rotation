package business

import (
	"math"
	"math/rand"
	"time"

	"github.com/pls87/creative-rotation/internal/storage/models"
)

const mistake = 1e-09

func NextCreative(stats []models.Stats) (creativeID models.ID) {
	rand.Seed(time.Now().Unix())
	zeroImp, totalImp := aggregate(stats)
	if len(zeroImp) > 0 {
		return zeroImp[rand.Intn(len(zeroImp))].CreativeID //nolint: gosec
	}
	var cur float64
	max := math.Inf(-1)
	creatives := make([]models.ID, 0, 10)
	for _, s := range stats {
		cur = valueToMaximize(s, totalImp)
		if max-cur > mistake {
			continue
		}
		if cur-max > mistake {
			max = cur
			creatives = creatives[:0]
		}
		creatives = append(creatives, s.CreativeID)
	}
	return creatives[rand.Intn(len(creatives))] //nolint: gosec
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
