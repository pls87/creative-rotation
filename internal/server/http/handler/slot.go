package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type SlotService struct {
	logger  *logger.Logger
	slotApp app.SlotApplication
	resp    *response
	helper  *helpers
}

func (s *SlotService) All(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slots, err := s.slotApp.All(ctx)
	if err != nil {
		s.resp.internalServerError(ctx, w, "Unexpected error while getting slots from storage", err)
		return
	}

	s.resp.json(ctx, w, map[string][]models.Slot{"slots": slots})
}

func (s *SlotService) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var toCreate models.Slot
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		s.resp.badRequest(ctx, w, "failed to parse slot body", err)
		return
	}

	if toCreate.Desc == "" {
		s.resp.badRequest(ctx, w, "slot description can't be empty", err)
		return
	}

	created, err := s.slotApp.New(ctx, toCreate)
	if err != nil {
		s.resp.internalServerError(ctx, w, "Unexpected error while saving slot to storage", err)
		return
	}

	s.resp.json(ctx, w, map[string]models.Slot{"slot": created})
}

func (s *SlotService) Creatives(w http.ResponseWriter, r *http.Request) {
	var id models.ID
	var ok bool
	if id, ok = s.helper.handleURLParamID(w, r, "slot_id"); !ok {
		return
	}
	ctx := r.Context()
	creatives, err := s.slotApp.Creatives(ctx, id)
	if err != nil {
		s.resp.internalServerError(ctx, w,
			fmt.Sprintf("Unexpected error while getting creatives for slot '%d'", id), err)
		return
	}

	s.resp.json(ctx, w, map[string][]models.Creative{"creatives": creatives})
}
