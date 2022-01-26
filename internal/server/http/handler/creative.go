package handler

import (
	"net/http"

	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/logger"
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

	s.resp.json(ctx, w, creatives)
}
