package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

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

func (s *CreativeService) Slots(w http.ResponseWriter, r *http.Request) {
	var id models.ID
	var ok bool
	if id, ok = s.handleURLParamID(w, r, "creative_id"); !ok {
		return
	}
	ctx := r.Context()
	slots, err := s.creativeApp.Slots(ctx, id)
	if err != nil {
		s.resp.internalServerError(ctx, w,
			fmt.Sprintf("Unexpected error while getting slots for creative '%d'", id), err)
		return
	}

	s.resp.json(ctx, w, map[string][]models.Slot{"slots": slots})
}

func (s *CreativeService) AllCreativeSlots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slotCreatives, err := s.creativeApp.AllCreativeSlots(ctx)
	if err != nil {
		s.resp.internalServerError(ctx, w, "Unexpected error while getting creative-slots '%d'", err)
		return
	}

	s.resp.json(ctx, w, map[string][]models.SlotCreative{"slot_creatives": slotCreatives})
}

func (s *CreativeService) handleURLParamID(w http.ResponseWriter,
	r *http.Request, name string) (id models.ID, ok bool) {
	vars := mux.Vars(r)

	tempID, e := strconv.Atoi(vars[name])
	if e != nil || tempID <= 0 {
		s.resp.badRequest(r.Context(), w, fmt.Sprintf("malformed %s", name), e)
		return 0, false
	}

	return models.ID(tempID), true
}

func (s *CreativeService) handleSlotBody(w http.ResponseWriter, r *http.Request) (slotID models.ID, ok bool) {
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
	if creativeID, ok = s.handleURLParamID(w, r, "creative_id"); !ok {
		return
	}
	if slotID, ok = s.handleSlotBody(w, r); !ok {
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
	if creativeID, ok = s.handleURLParamID(w, r, "creative_id"); !ok {
		return
	}
	if slotID, ok = s.handleURLParamID(w, r, "slot_id"); !ok {
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

func (s *CreativeService) TrackImpression(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var toCreate models.Impression
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		s.resp.badRequest(ctx, w, "failed to parse impression body", err)
		return
	}

	err = s.creativeApp.TrackImpression(ctx, toCreate)
	if err != nil {
		s.resp.internalServerError(ctx, w, "Unexpected error while saving impression to storage", err)
		return
	}

	s.resp.json(ctx, w, map[string]models.Impression{"impression": toCreate})
}

func (s *CreativeService) Next(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slotID, ok := s.handleIDQuery(ctx, "slot_id", w, r)
	if !ok {
		return
	}
	segmentID, ok := s.handleIDQuery(ctx, "segment_id", w, r)
	if !ok {
		return
	}

	creative, err := s.creativeApp.Next(ctx, slotID, segmentID)
	if err != nil || creative.ID <= 0 {
		s.resp.internalServerError(ctx, w, "Unexpected error while getting next creative", err)
		return
	}

	s.resp.json(ctx, w, map[string]models.Creative{"creative": creative})
}

func (s *CreativeService) handleIDQuery(ctx context.Context, param string, w http.ResponseWriter,
	r *http.Request) (models.ID, bool) {
	IDs := r.URL.Query()[param]
	if len(IDs) == 0 {
		s.resp.badRequest(ctx, w, fmt.Sprintf("%s isn't specified", param), nil)
		return 0, false
	}
	if len(IDs) != 1 {
		s.resp.badRequest(ctx, w, fmt.Sprintf("more than one %s were passed", param), nil)
		return 0, false
	}

	ID, err := strconv.Atoi(IDs[0])
	if err != nil || ID <= 0 {
		s.resp.badRequest(ctx, w, fmt.Sprintf("malformed %s", param), err)
		return 0, false
	}

	return models.ID(ID), true
}
