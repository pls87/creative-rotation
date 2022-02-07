package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/server/http/handler/helpers"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type CreativeService struct {
	logger *logger.Logger
	app    app.CreativeApplication
	resp   *helpers.Responder
	helper *helpers.ParamHelper
}

func (s *CreativeService) All(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	creatives, err := s.app.All(ctx)
	if err != nil {
		s.resp.InternalServerError(ctx, w, "Unexpected error while getting creatives from storage", err)
		return
	}

	s.resp.JSON(ctx, w, map[string][]models.Creative{"creatives": creatives})
}

func (s *CreativeService) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var toCreate models.Creative
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		s.resp.BadRequest(ctx, w, "failed to parse creative body", err)
		return
	}

	if toCreate.Desc == "" {
		s.resp.BadRequest(ctx, w, "creative description can't be empty", err)
		return
	}

	created, err := s.app.New(ctx, toCreate)
	if err != nil {
		s.resp.InternalServerError(ctx, w, "Unexpected error while saving creative to storage", err)
		return
	}

	s.resp.JSON(ctx, w, created)
}

func (s *CreativeService) Slots(w http.ResponseWriter, r *http.Request) {
	var id models.ID
	var ok bool
	if id, ok = s.helper.HandleURLParamID(w, r, "creative_id"); !ok {
		return
	}
	ctx := r.Context()
	slots, err := s.app.Slots(ctx, id)
	if err != nil {
		s.resp.InternalServerError(ctx, w,
			fmt.Sprintf("Unexpected error while getting slots for creative '%d'", id), err)
		return
	}

	s.resp.JSON(ctx, w, map[string][]models.Slot{"slots": slots})
}

func (s *CreativeService) AllCreativeSlots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slotCreatives, err := s.app.AllCreativeSlots(ctx)
	if err != nil {
		s.resp.InternalServerError(ctx, w, "Unexpected error while getting creative-slots '%d'", err)
		return
	}

	s.resp.JSON(ctx, w, map[string][]models.SlotCreative{"slot_creatives": slotCreatives})
}

func (s *CreativeService) handleSlotBody(w http.ResponseWriter, r *http.Request) (slotID models.ID, ok bool) {
	var slot models.Slot
	err := json.NewDecoder(r.Body).Decode(&slot)
	if err != nil || slot.ID <= 0 {
		s.resp.BadRequest(r.Context(), w, "failed to parse slot", err)
		return 0, false
	}

	return slot.ID, true
}

func (s *CreativeService) AddToSlot(w http.ResponseWriter, r *http.Request) {
	var creativeID, slotID models.ID
	var ok bool
	if creativeID, ok = s.helper.HandleURLParamID(w, r, "creative_id"); !ok {
		return
	}
	if slotID, ok = s.handleSlotBody(w, r); !ok {
		return
	}

	ctx := r.Context()
	err := s.app.AddToSlot(ctx, creativeID, slotID)
	if err != nil {
		s.resp.InternalServerError(ctx, w, "Unexpected error while adding creative to slot", err)
		return
	}

	res := models.SlotCreative{SlotID: slotID, CreativeID: creativeID}
	s.resp.JSON(ctx, w, res)
}

func (s *CreativeService) RemoveFromSlot(w http.ResponseWriter, r *http.Request) {
	var creativeID, slotID models.ID
	var ok bool
	if creativeID, ok = s.helper.HandleURLParamID(w, r, "creative_id"); !ok {
		return
	}
	if slotID, ok = s.helper.HandleURLParamID(w, r, "slot_id"); !ok {
		return
	}

	ctx := r.Context()
	err := s.app.RemoveFromSlot(ctx, creativeID, slotID)
	if err != nil {
		s.resp.InternalServerError(ctx, w, "Unexpected error while removing creative from slot", err)
		return
	}

	s.resp.JSON(ctx, w, true)
}

func (s *CreativeService) TrackConversion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var toCreate models.Conversion
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		s.resp.BadRequest(ctx, w, "failed to parse conversion body", err)
		return
	}

	err = s.app.TrackConversion(ctx, toCreate)
	if err != nil {
		s.resp.InternalServerError(ctx, w, "Unexpected error while saving conversion to storage", err)
		return
	}

	s.resp.JSON(ctx, w, toCreate)
}

func (s *CreativeService) TrackImpression(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var toCreate models.Impression
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		s.resp.BadRequest(ctx, w, "failed to parse impression body", err)
		return
	}

	err = s.app.TrackImpression(ctx, toCreate)
	if err != nil {
		s.resp.InternalServerError(ctx, w, "Unexpected error while saving impression to storage", err)
		return
	}

	s.resp.JSON(ctx, w, toCreate)
}

func (s *CreativeService) Next(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slotID, ok := s.helper.HandleIDQuery(ctx, "slot_id", w, r)
	if !ok {
		return
	}
	segmentID, ok := s.helper.HandleIDQuery(ctx, "segment_id", w, r)
	if !ok {
		return
	}

	creative, err := s.app.Next(ctx, slotID, segmentID)
	if err != nil || creative.ID <= 0 {
		s.resp.InternalServerError(ctx, w, "Unexpected error while getting next creative", err)
		return
	}

	s.resp.JSON(ctx, w, creative)
}
