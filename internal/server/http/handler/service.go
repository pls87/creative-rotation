package handler

import (
	"net/http"

	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/logger"
)

type Service struct {
	creatives *CreativeService
	slots     *SlotService
	segments  *SegmentService
	logger    *logger.Logger
	resp      *response
}

func NewService(app app.Application, logger *logger.Logger) *Service {
	resp := &response{logger: logger}
	return &Service{
		creatives: &CreativeService{logger: logger, creativeApp: app.Creatives(), resp: resp},
		slots:     &SlotService{logger: logger, slotApp: app.Slots(), resp: resp},
		segments:  &SegmentService{logger: logger, segmentApp: app.Segments(), resp: resp},
		logger:    logger,
		resp:      resp,
	}
}

func (s *Service) Creatives() *CreativeService {
	return s.creatives
}

func (s *Service) Slots() *SlotService {
	return s.slots
}

func (s *Service) Segments() *SegmentService {
	return s.segments
}

func (s *Service) Noop(w http.ResponseWriter, r *http.Request) {
	s.resp.text(r.Context(), w, "It Works!")
}
