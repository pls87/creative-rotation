package handler

import (
	"net/http"

	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/logger"
)

type CreativeService struct {
	logger      *logger.Logger
	creativeApp app.CreativeApplication
}

func (s *CreativeService) All(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("All creatives!")); err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		s.logger.Errorf("Couldn't write an HTTP response: %s", err)
	}
}
