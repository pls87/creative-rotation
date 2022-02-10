//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"testing"
	"time"

	"github.com/pls87/creative-rotation/internal/tests/integration/helpers"
	"github.com/stretchr/testify/suite"
)

type probability struct {
	n  int
	of int
}

var (
	overallImpressionRate = probability{n: 4, of: 5}
	conversionRates       = map[int][]probability{
		1: { // lego technic
			probability{n: 5, of: 100}, // Girl
			probability{n: 2, of: 10},  // Boy
			probability{n: 1, of: 10},  // Man
			probability{n: 2, of: 100}, // Woman
		},
		2: { // lego friends
			probability{n: 2, of: 10},  // Girl
			probability{n: 5, of: 100}, // Boy
			probability{n: 2, of: 100}, // Man
			probability{n: 3, of: 100}, // Woman
		},
		3: { // kia soul
			probability{n: 2, of: 100}, // Girl
			probability{n: 5, of: 10},  // Boy
			probability{n: 1, of: 10},  // Man
			probability{n: 3, of: 10},  // Woman
		},
		4: { // chevrolet tahoe
			probability{n: 2, of: 100}, // Girl
			probability{n: 5, of: 100}, // Boy
			probability{n: 3, of: 10},  // Man
			probability{n: 1, of: 10},  // Woman
		},
		5: { // chanel chance
			probability{n: 5, of: 100}, // Girl
			probability{n: 1, of: 100}, // Boy
			probability{n: 1, of: 10},  // Man
			probability{n: 3, of: 10},  // Woman
		},
		6: { // dior homme
			probability{n: 1, of: 100}, // Girl
			probability{n: 2, of: 100}, // Boy
			probability{n: 3, of: 10},  // Man
			probability{n: 1, of: 10},  // Woman
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
			probability{n: 10, of: 100}, // Girl
			probability{n: 10, of: 100}, // Boy
			probability{n: 40, of: 100}, // Man
			probability{n: 40, of: 100}, // Woman
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
	return rand.Intn(4)
}

func selectSegment(slot int) int {
	ar := auditoryRates[slot]
	r := rand.Intn(100)
	cur := 0
	for i, v := range ar {
		if r >= cur && r < cur+v.n {
			return i
		}
		cur = cur + v.n
	}
	fmt.Println("slot: ", slot, "rand: ", r)
	return -1
}

type UCB1Suite struct {
	BaseSuite
	statsH    *helpers.StatHelper
	creativeH *helpers.CreativeHelper
}

func (s *UCB1Suite) run() {
	start := time.Now()
	numOfImpressions := 0
	for time.Since(start) < 10*time.Second {
		slot := selectSlot()
		segment := selectSegment(slot + 1)
		s.GreaterOrEqual(segment, 0)
		creative := s.nextCreative(segment+1, slot+1)

		if rand.Intn(overallImpressionRate.of) >= overallImpressionRate.n {
			continue
		}

		time.Sleep(10 * time.Millisecond)

		s.trackImpression(creative, slot+1, segment+1)
		numOfImpressions++
		convProb := conversionRates[creative][segment]
		if rand.Intn(convProb.of) >= convProb.n {
			continue
		}
		time.Sleep(10 * time.Millisecond)
		s.trackConversion(creative, slot+1, segment+1)
	}
	fmt.Println(numOfImpressions)
}

func (s *UCB1Suite) TestBusiness() {
	s.run()
}

func (s *UCB1Suite) nextCreative(segment, slot int) int {
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
