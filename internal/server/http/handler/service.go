package handler

import (
	"net/http"

	"github.com/pls87/creative-rotation/internal/app"

	"github.com/pls87/creative-rotation/internal/logger"
)

type Service struct {
	creatives *CreativeService
	logger    *logger.Logger
}

func NewService(app app.Application, logger *logger.Logger) *Service {
	return &Service{
		creatives: &CreativeService{
			logger:      logger,
			creativeApp: app.Creatives(),
		},
		logger: logger,
	}
}

func (s *Service) Creatives() *CreativeService {
	return s.creatives
}

func (s *Service) Noop(w http.ResponseWriter, r *http.Request) {
	if _, err := w.Write([]byte("It works!")); err == nil {
		w.WriteHeader(http.StatusOK)
	} else {
		s.logger.Errorf("Couldn't write an HTTP response: %s", err)
	}
}
