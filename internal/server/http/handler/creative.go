package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type CreativeService struct {
	logger      *logger.Logger
	creativeApp app.CreativeApplication
	resp        *response
}

func (s *CreativeService) All(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creatives, err := s.creativeApp.All(ctx)
	if err != nil {
		s.resp.internalServerError(ctx, w, "Unexpected error while getting creatives from storage", err)
		return
	}

	s.resp.json(ctx, w, map[string][]models.Creative{"creatives": creatives})
}

func (s *CreativeService) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var toCreate models.Creative
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		s.resp.badRequest(ctx, w, "failed to parse creative body", err)
		return
	}

	created, err := s.creativeApp.New(ctx, toCreate)
	if err != nil {
		s.resp.internalServerError(ctx, w, "Unexpected error while saving creative to storage", err)
		return
	}

	s.resp.json(ctx, w, map[string]models.Creative{"creative": created})
}

func (s *CreativeService) handleCreativeParams(w http.ResponseWriter, r *http.Request) (id models.ID, ok bool) {
	vars := mux.Vars(r)

	creativeID, e := strconv.Atoi(vars["id"])
	if e != nil || creativeID <= 0 {
		s.resp.badRequest(r.Context(), w, "malformed creative id", e)
		return 0, false
	}

	return models.ID(creativeID), true
}

func (s *CreativeService) handleSlotParams(w http.ResponseWriter, r *http.Request) (slotID models.ID, ok bool) {
	var slot models.Slot
	err := json.NewDecoder(r.Body).Decode(&slot)
	if err != nil || slot.ID <= 0 {
		s.resp.badRequest(r.Context(), w, "failed to parse slot", err)
		return 0, false
	}

	return slot.ID, true
}

func (s *CreativeService) AddToSlot(w http.ResponseWriter, r *http.Request) {
	var creativeID, slotID models.ID
	var ok bool
	if creativeID, ok = s.handleCreativeParams(w, r); !ok {
		return
	}
	if slotID, ok = s.handleSlotParams(w, r); !ok {
		return
	}

	ctx := r.Context()
	err := s.creativeApp.AddToSlot(ctx, creativeID, slotID)
	if err != nil {
		s.resp.internalServerError(ctx, w, "Unexpected error while adding creative to slot", err)
		return
	}

	res := models.SlotCreative{SlotID: slotID, CreativeID: creativeID}
	s.resp.json(ctx, w, res)
}

func (s *CreativeService) RemoveFromSlot(w http.ResponseWriter, r *http.Request) {
	var creativeID, slotID models.ID
	var ok bool
	if creativeID, ok = s.handleCreativeParams(w, r); !ok {
		return
	}
	if slotID, ok = s.handleSlotParams(w, r); !ok {
		return
	}

	ctx := r.Context()
	err := s.creativeApp.RemoveFromSlot(ctx, creativeID, slotID)
	if err != nil {
		s.resp.internalServerError(ctx, w, "Unexpected error while removing creative from slot", err)
		return
	}

	s.resp.json(ctx, w, true)
}

func (s *CreativeService) TrackConversion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var toCreate models.Conversion
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		s.resp.badRequest(ctx, w, "failed to parse conversion body", err)
		return
	}

	err = s.creativeApp.TrackConversion(ctx, toCreate)
	if err != nil {
		s.resp.internalServerError(ctx, w, "Unexpected error while saving conversion to storage", err)
		return
	}

	s.resp.json(ctx, w, map[string]models.Conversion{"conversion": toCreate})
}

func (s *CreativeService) Next(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	slotID, err := strconv.Atoi(strings.Join(r.URL.Query()["slot_id"], ""))
	if err != nil || slotID <= 0 {
		s.resp.badRequest(ctx, w, "slot id isn't specified", err)
		return
	}

	segmentID, err := strconv.Atoi(strings.Join(r.URL.Query()["segment_id"], ""))
	if err != nil || segmentID <= 0 {
		s.resp.badRequest(ctx, w, "segment id isn't specified", err)
		return
	}

	creative, err := s.creativeApp.Next(ctx, models.ID(slotID), models.ID(segmentID))
	if err != nil || creative.ID <= 0 {
		s.resp.internalServerError(ctx, w, "Unexpected error while getting next creative", err)
		return
	}

	s.resp.json(ctx, w, map[string]models.Creative{"creative": creative})
}
