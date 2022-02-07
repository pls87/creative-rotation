package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/server/http/handler/helpers"
	"github.com/pls87/creative-rotation/internal/storage/basic"
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
		s.resp.InternalServerError(ctx, w, helpers.UnexpectedErrorGetCreatives, err)
		return
	}

	s.resp.JSON(ctx, w, helpers.CreativeCollection{Creatives: creatives})
}

func (s *CreativeService) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var toCreate models.Creative
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		s.resp.BadRequest(ctx, w, helpers.BadRequestFailedParseCreative, err)
		return
	}

	if toCreate.Desc == "" {
		s.resp.BadRequest(ctx, w, helpers.BadRequestFailedEmptyDescCreative, err)
		return
	}

	created, err := s.app.New(ctx, toCreate)
	if err != nil {
		s.resp.InternalServerError(ctx, w, helpers.UnexpectedErrorSavingCreative, err)
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
		s.resp.InternalServerError(ctx, w, helpers.UnexpectedErrorGetSlots, err)
		return
	}

	s.resp.JSON(ctx, w, helpers.SlotCollection{Slots: slots})
}

func (s *CreativeService) AllCreativeSlots(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slotCreatives, err := s.app.AllSlotCreatives(ctx)
	if err != nil {
		s.resp.InternalServerError(ctx, w, helpers.UnexpectedErrorGetSlotCreatives, err)
		return
	}

	s.resp.JSON(ctx, w, helpers.SCCollection{SC: slotCreatives})
}

func (s *CreativeService) handleSlotBody(w http.ResponseWriter, r *http.Request) (slotID models.ID, ok bool) {
	var slot models.Slot
	err := json.NewDecoder(r.Body).Decode(&slot)
	if err != nil || slot.ID <= 0 {
		s.resp.BadRequest(r.Context(), w, helpers.BadRequestFailedParseSlot, err)
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

	switch {
	case err == nil:
		s.resp.JSON(ctx, w, models.SlotCreative{SlotID: slotID, CreativeID: creativeID})
	case errors.Is(err, basic.ErrCreativeAlreadyInSlot):
		s.resp.Conflict(ctx, w, helpers.ConflictCreativeAlreadyInSlot, err)
	default:
		s.resp.InternalServerError(ctx, w, helpers.UnexpectedErrorSavingSlotCreative, err)
	}
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
	switch {
	case err == nil:
		s.resp.JSON(ctx, w, true)
	case errors.Is(err, basic.ErrCreativeNotInSlot):
		s.resp.NotFound(ctx, w, helpers.NotFoundCreativeNotInSlot, err)
	default:
		s.resp.InternalServerError(ctx, w, helpers.UnexpectedErrorDeleteCreativeFromSlot, err)
	}
}

func (s *CreativeService) TrackConversion(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var toCreate models.Conversion
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		s.resp.BadRequest(ctx, w, helpers.BadRequestFailedParseConversion, err)
		return
	}

	err = s.app.TrackConversion(ctx, toCreate)
	if err != nil {
		s.resp.InternalServerError(ctx, w, helpers.UnexpectedErrorSavingConversion, err)
		return
	}

	s.resp.JSON(ctx, w, toCreate)
}

func (s *CreativeService) TrackImpression(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var toCreate models.Impression
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		s.resp.BadRequest(ctx, w, helpers.BadRequestFailedParseImpression, err)
		return
	}

	err = s.app.TrackImpression(ctx, toCreate)
	if err != nil {
		s.resp.InternalServerError(ctx, w, helpers.UnexpectedErrorSavingImpression, err)
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
		s.resp.InternalServerError(ctx, w, helpers.UnexpectedErrorGetNextCreative, err)
		return
	}

	s.resp.JSON(ctx, w, creative)
}
