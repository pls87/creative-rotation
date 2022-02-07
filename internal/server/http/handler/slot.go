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

type SlotService struct {
	logger *logger.Logger
	app    app.SlotApplication
	resp   *helpers.Responder
	helper *helpers.ParamHelper
}

func (s *SlotService) All(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slots, err := s.app.All(ctx)
	if err != nil {
		s.resp.InternalServerError(ctx, w, "Unexpected error while getting slots from storage", err)
		return
	}

	s.resp.JSON(ctx, w, map[string][]models.Slot{"slots": slots})
}

func (s *SlotService) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var toCreate models.Slot
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		s.resp.BadRequest(ctx, w, "failed to parse slot body", err)
		return
	}

	if toCreate.Desc == "" {
		s.resp.BadRequest(ctx, w, "slot description can't be empty", err)
		return
	}

	created, err := s.app.New(ctx, toCreate)
	if err != nil {
		s.resp.InternalServerError(ctx, w, "Unexpected error while saving slot to storage", err)
		return
	}

	s.resp.JSON(ctx, w, created)
}

func (s *SlotService) Creatives(w http.ResponseWriter, r *http.Request) {
	var id models.ID
	var ok bool
	if id, ok = s.helper.HandleURLParamID(w, r, "slot_id"); !ok {
		return
	}
	ctx := r.Context()
	creatives, err := s.app.Creatives(ctx, id)
	if err != nil {
		s.resp.InternalServerError(ctx, w,
			fmt.Sprintf("Unexpected error while getting creatives for slot '%d'", id), err)
		return
	}

	s.resp.JSON(ctx, w, map[string][]models.Creative{"creatives": creatives})
}
