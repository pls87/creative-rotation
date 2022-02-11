//go:build integration
// +build integration

// this test shouldn't be considered as always green - it uses random functions and probabilities
// if it fails often it may mean that algorithm doesn't work fine
// if you need to introduce this into CI than it should be run several times and allow fail once or twice
package integration

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/pls87/creative-rotation/internal/tests/integration/helpers"
	"github.com/stretchr/testify/suite"
)

var (
	sleepGap = 5 * time.Millisecond
	testTime = 20 * time.Second
)

type probability struct {
	n  int
	of int
}

var (
	overallImpressionRate = probability{n: 4, of: 5}
	// CTRs are much higher than in reality but OK for tests purposes
	conversionRates = map[int][]probability{
		1: { // lego technic
			probability{n: 30, of: 100}, // Girl
			probability{n: 80, of: 100}, // Boy
			probability{n: 50, of: 100}, // Man
			probability{n: 10, of: 100}, // Woman
		},
		2: { // lego friends
			probability{n: 60, of: 100}, // Girl
			probability{n: 25, of: 100}, // Boy
			probability{n: 5, of: 100},  // Man
			probability{n: 45, of: 100}, // Woman
		},
		3: { // kia soul
			probability{n: 5, of: 100},  // Girl
			probability{n: 10, of: 100}, // Boy
			probability{n: 20, of: 100}, // Man
			probability{n: 40, of: 100}, // Woman
		},
		4: { // chevrolet tahoe
			probability{n: 5, of: 100},  // Girl
			probability{n: 10, of: 100}, // Boy
			probability{n: 40, of: 100}, // Man
			probability{n: 20, of: 100}, // Woman
		},
		5: { // chanel chance
			probability{n: 10, of: 100}, // Girl
			probability{n: 5, of: 100},  // Boy
			probability{n: 20, of: 100}, // Man
			probability{n: 40, of: 100}, // Woman
		},
		6: { // dior homme
			probability{n: 5, of: 100},  // Girl
			probability{n: 10, of: 100}, // Boy
			probability{n: 40, of: 100}, // Man
			probability{n: 20, of: 100}, // Woman
		},
	}
	auditoryRates = map[int][]probability{
		1: { // drom.ru
			probability{n: 2, of: 100},  // Girl
			probability{n: 8, of: 100},  // Boy
			probability{n: 60, of: 100}, // Man
			probability{n: 30, of: 100}, // Woman
		},
		2: { // ozon.ru
			probability{n: 25, of: 100}, // Girl
			probability{n: 25, of: 100}, // Boy
			probability{n: 25, of: 100}, // Man
			probability{n: 25, of: 100}, // Woman
		},
		3: { // toys.ru
			probability{n: 30, of: 100}, // Girl
			probability{n: 30, of: 100}, // Boy
			probability{n: 20, of: 100}, // Man
			probability{n: 20, of: 100}, // Woman
		},
		4: { // letu.ru
			probability{n: 8, of: 100},  // Girl
			probability{n: 2, of: 100},  // Boy
			probability{n: 30, of: 100}, // Man
			probability{n: 60, of: 100}, // Woman
		},
	}
)

func selectSlot() int {
	return getRand(4)
}

func selectSegment(slot int) int {
	ar := auditoryRates[slot]
	r := getRand(100)
	cur := 0
	for i, v := range ar {
		if r >= cur && r < cur+v.n {
			return i
		}
		cur = cur + v.n
	}
	return -1
}

func conversionRate(stats helpers.Stats) float64 {
	return float64(stats.Conversions) / float64(stats.Impressions)
}

func getRand(n int) int {
	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(n)))
	if err != nil {
		panic(err)
	}
	return int(nBig.Int64())
}

type UCB1Suite struct {
	BaseSuite
	statsH    *helpers.StatHelper
	creativeH *helpers.CreativeHelper
}

func (s *UCB1Suite) run(imp, conv *int64) {
	start := time.Now()
	for time.Since(start) < testTime {
		slotIndex := selectSlot()
		slotID := slotIndex + 1
		segmentIndex := selectSegment(slotID)
		s.GreaterOrEqual(segmentIndex, 0)
		segmentID := segmentIndex + 1

		creativeID := s.nextCreative(slotID, segmentID)

		if getRand(overallImpressionRate.of) >= overallImpressionRate.n {
			continue
		}

		time.Sleep(sleepGap)

		s.trackImpression(creativeID, slotID, segmentID)
		atomic.AddInt64(imp, 1)
		convProb := conversionRates[creativeID][segmentIndex]
		if getRand(convProb.of) >= convProb.n {
			continue
		}
		time.Sleep(sleepGap)
		s.trackConversion(creativeID, slotID, segmentID)
		atomic.AddInt64(conv, 1)
		time.Sleep(sleepGap)
	}
}

func (s *UCB1Suite) TestBusiness() {
	wg := sync.WaitGroup{}
	var imp, conv int64
	wg.Add(1)
	go func() {
		s.run(&imp, &conv)
		wg.Done()
	}()
	wg.Add(1)
	go func() {
		s.run(&imp, &conv)
		wg.Done()
	}()
	wg.Wait()

	time.Sleep(1 * time.Second)

	stats := s.getStats("")
	var i, c int
	for _, v := range stats {
		i += v.Impressions
		c += v.Conversions
	}
	s.Equalf(imp, int64(i), "%d impressions was done but %d was tracked", imp, i)
	s.Equalf(conv, int64(c), "%d conversions was done but %d was tracked", conv, c)

	s.Empty(s.getStats("WHERE impressions=0"))

	stats = s.getStats("WHERE slot_id=2 AND creative_id=1") // ozon & lego technic
	var boyStat, girlStat, manStat, womanStat helpers.Stats
	for _, v := range stats {
		switch v.SegmentID {
		case 1:
			girlStat = v
		case 2:
			boyStat = v
		case 3:
			manStat = v
		case 4:
			womanStat = v
		}
	}
	s.Greater(conversionRate(boyStat), conversionRate(girlStat), "boy stat should be better than girls")
	s.Greater(conversionRate(boyStat), conversionRate(manStat), "boy stat should be better than man")
	s.Greater(conversionRate(manStat), conversionRate(womanStat), "man stat should be better than woman")
}

func (s *UCB1Suite) getStats(where string) (stats []helpers.Stats) {
	err := s.client.DB.Select(&stats, fmt.Sprintf(`SELECT creative_id, slot_id, segment_id,
		CASE WHEN impressions is NULL THEN 0 ELSE impressions END as impressions,
		CASE WHEN conversions is NULL THEN 0 ELSE conversions END as conversions
		FROM "stats" %s`, where))
	s.NoError(err)
	return stats
}

func (s *UCB1Suite) nextCreative(slot, segment int) int {
	code, body, err := s.creativeH.Next(slot, segment)
	s.NoError(err)
	s.Equal(code, http.StatusOK)
	creative := helpers.Entity{}
	s.NoError(json.Unmarshal(body, &creative))

	return creative.ID
}

func (s *UCB1Suite) trackImpression(creativeID, slotID, segmentID int) {
	code, _, err := s.statsH.TrackImpression(creativeID, slotID, segmentID)
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equalf(http.StatusOK, code, "status %d expected but was %d", http.StatusOK, code)
}

func (s *UCB1Suite) trackConversion(creativeID, slotID, segmentID int) {
	code, _, err := s.statsH.TrackConversion(creativeID, slotID, segmentID)
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equalf(http.StatusOK, code, "status %d expected but was %d", http.StatusOK, code)
}

func (s *UCB1Suite) SetupSuite() {
	s.BaseSuite.SetupSuite()
	s.creativeH = helpers.NewCreativeHelper(s.baseURL)
	s.statsH = helpers.NewStatHelper(s.baseURL)
	s.NoError(s.client.RunFile("./sql/seed.sql"))
}

func TestBusinessUCB1(t *testing.T) {
	suite.Run(t, new(UCB1Suite))
}
