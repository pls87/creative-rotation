package handler

import (
	"encoding/json"
	"net/http"

	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/server/http/handler/helpers"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type SegmentService struct {
	logger *logger.Logger
	app    app.SegmentApplication
	resp   *helpers.Responder
	helper *helpers.ParamHelper
}

func (s *SegmentService) All(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	segments, err := s.app.All(ctx)
	if err != nil {
		s.resp.InternalServerError(ctx, w, helpers.UnexpectedErrorGetSegments, err)
		return
	}

	s.resp.JSON(ctx, w, helpers.SegmentCollection{Segments: segments})
}

func (s *SegmentService) New(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var toCreate models.Segment
	err := json.NewDecoder(r.Body).Decode(&toCreate)
	if err != nil {
		s.resp.BadRequest(ctx, w, helpers.BadRequestFailedParseSegment, err)
		return
	}

	if toCreate.Desc == "" {
		s.resp.BadRequest(ctx, w, helpers.BadRequestFailedEmptyDescSegment, err)
		return
	}

	created, err := s.app.New(ctx, toCreate)
	if err != nil {
		s.resp.InternalServerError(ctx, w, helpers.UnexpectedErrorSavingSegment, err)
		return
	}

	s.resp.JSON(ctx, w, created)
}
