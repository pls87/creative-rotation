package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type SegmentService struct {
	logger     *logger.Logger
	segmentApp app.SegmentApplication
	resp       *response
}

func (s *SegmentService) All(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	segments, err := s.segmentApp.All(ctx)
	if err != nil {
		s.resp.internalServerError(ctx, w, "Unexpected error while getting segments from storage", err)
		return
	}

	s.resp.json(ctx, w, map[string][]models.Segment{"segments": segments})
}

func (s *SegmentService) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var toCreate models.Segment
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		s.resp.badRequest(ctx, w, "failed to parse segment body", err)
		return
	}

	if toCreate.Desc == "" {
		s.resp.badRequest(ctx, w, "segment description can't be empty", err)
		return
	}

	created, err := s.segmentApp.New(ctx, toCreate)
	if err != nil {
		s.resp.internalServerError(ctx, w, "Unexpected error while saving segment to storage", err)
		return
	}

	s.resp.json(ctx, w, map[string]models.Segment{"segment": created})
}
