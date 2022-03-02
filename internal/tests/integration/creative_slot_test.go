//go:build integration
// +build integration

package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/pls87/creative-rotation/internal/tests/integration/helpers"
	"github.com/stretchr/testify/suite"
)

type CreativeSlotSuite struct {
	BaseSuite
	entitiesH  *helpers.EntityHelper
	creativesH *helpers.CreativeHelper
	slotsH     *helpers.SlotHelper

	creatives []helpers.Entity
	slots     []helpers.Entity
}

func (s *CreativeSlotSuite) SetupSuite() {
	s.BaseSuite.SetupSuite()
	s.entitiesH = helpers.NewEntityHelper(s.baseURL)
	s.creativesH = helpers.NewCreativeHelper(s.baseURL)
	s.slotsH = helpers.NewSlotHelper(s.baseURL)
}

func (s *CreativeSlotSuite) SetupTest() {
	s.creatives = make([]helpers.Entity, 0, 2)
	s.slots = make([]helpers.Entity, 0, 2)
	s.creatives = append(s.creatives, s.entitiesH.New(s.T(), "creative", gofakeit.BuzzWord()))
	s.creatives = append(s.creatives, s.entitiesH.New(s.T(), "creative", gofakeit.BuzzWord()))
	s.slots = append(s.slots, s.entitiesH.New(s.T(), "slot", gofakeit.BuzzWord()))
	s.slots = append(s.slots, s.entitiesH.New(s.T(), "slot", gofakeit.BuzzWord()))
}

func (s *CreativeSlotSuite) TestAddToSlot() {
	s.addToSlot(s.creatives[0].ID, s.slots[1].ID)
	s.addToSlot(s.creatives[1].ID, s.slots[0].ID)
	s.addToSlot(s.creatives[1].ID, s.slots[1].ID)

	code, body, err := s.creativesH.AllCreativeSlots()
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equalf(http.StatusOK, code, "status %d expected but was %d", http.StatusOK, code)
	var target map[string][]helpers.SlotCreative
	s.NoError(json.Unmarshal(body, &target))
	slotCreatives, ok := target["slot_creatives"]
	s.Truef(ok, "slot_creatives key wasn't found in the response")
	var f1, f2, f3 bool
	for _, v := range slotCreatives {
		if v.SlotID == s.slots[1].ID && v.SlotDesc == s.slots[1].Desc &&
			v.CreativeID == s.creatives[0].ID && v.CreativeDesc == s.creatives[0].Desc {
			f1 = true
		}
		if v.SlotID == s.slots[0].ID && v.SlotDesc == s.slots[0].Desc &&
			v.CreativeID == s.creatives[1].ID && v.CreativeDesc == s.creatives[1].Desc {
			f2 = true
		}
		if v.SlotID == s.slots[1].ID && v.SlotDesc == s.slots[1].Desc &&
			v.CreativeID == s.creatives[1].ID && v.CreativeDesc == s.creatives[1].Desc {
			f3 = true
		}
	}

	s.Truef(f1, "slot-creative %v - %v wasn't found in the response", s.creatives[0], s.slots[1])
	s.Truef(f2, "slot-creative %v - %v wasn't found in the response", s.creatives[1], s.slots[0])
	s.Truef(f3, "slot-creative %v - %v wasn't found in the response", s.creatives[1], s.slots[1])
}

func (s *CreativeSlotSuite) TestAddDuplicates() {
	s.addToSlot(s.creatives[0].ID, s.slots[1].ID)
	code, _, err := s.creativesH.AddToSlot(s.creatives[0].ID, s.slots[1].ID)
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equalf(http.StatusConflict, code, "status %d expected but was %d", http.StatusConflict, code)
}

func (s *CreativeSlotSuite) TestDeleteFromSlot() {
	s.addToSlot(s.creatives[0].ID, s.slots[0].ID)
	s.addToSlot(s.creatives[0].ID, s.slots[1].ID)
	code, _, err := s.creativesH.AddToSlot(s.creatives[0].ID, s.slots[1].ID)
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equalf(http.StatusConflict, code, "status %d expected but was %d", http.StatusConflict, code)
	code, _, err = s.creativesH.RemoveFromSlot(s.creatives[0].ID, s.slots[1].ID)
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equalf(http.StatusOK, code, "status %d expected but was %d", http.StatusOK, code)
	code, _, err = s.creativesH.RemoveFromSlot(s.creatives[0].ID, s.slots[1].ID)
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equalf(http.StatusNotFound, code, "status %d expected but was %d", http.StatusNotFound, code)
}

func (s *CreativeSlotSuite) TestCreativeSlots() {
	s.addToSlot(s.creatives[0].ID, s.slots[1].ID)
	s.addToSlot(s.creatives[1].ID, s.slots[0].ID)
	s.addToSlot(s.creatives[1].ID, s.slots[1].ID)

	code, body, err := s.creativesH.Slots(s.creatives[1].ID)
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equalf(http.StatusOK, code, "status %d expected but was %d", http.StatusOK, code)

	slots, err := helpers.ParseMany("slots", body)
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equal(2, len(slots))
	s.Equal(s.slots, slots)
}

func (s *CreativeSlotSuite) TestSlotCreatives() {
	s.addToSlot(s.creatives[0].ID, s.slots[1].ID)
	s.addToSlot(s.creatives[1].ID, s.slots[0].ID)
	s.addToSlot(s.creatives[1].ID, s.slots[1].ID)

	code, body, err := s.slotsH.Creatives(s.slots[1].ID)
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equalf(http.StatusOK, code, "status %d expected but was %d", http.StatusOK, code)

	creatives, err := helpers.ParseMany("creatives", body)
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equal(2, len(creatives))
	s.Equal(s.creatives, creatives)
}

func (s *CreativeSlotSuite) addToSlot(creativeID, slotID int) {
	code, _, err := s.creativesH.AddToSlot(creativeID, slotID)
	s.NoErrorf(err, "no error expected, but was: %s", err)
	s.Equalf(http.StatusOK, code, "status %d expected but was %d", http.StatusOK, code)
}

func TestCreativeSlot(t *testing.T) {
	suite.Run(t, new(CreativeSlotSuite))
}
