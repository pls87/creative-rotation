package handler

import (
	"net/http"

	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/server/http/handler/helpers"
)

type Service struct {
	creatives *CreativeService
	slots     *SlotService
	segments  *SegmentService
	logger    *logger.Logger
	resp      *helpers.Responder
}

func NewService(app app.Application, logger *logger.Logger) *Service {
	resp := helpers.NewResponder(logger)
	helper := helpers.NewParamHelper(resp)
	return &Service{
		creatives: &CreativeService{logger: logger, app: app.Creatives(), resp: resp, helper: helper},
		slots:     &SlotService{logger: logger, app: app.Slots(), resp: resp, helper: helper},
		segments:  &SegmentService{logger: logger, app: app.Segments(), resp: resp, helper: helper},
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
	s.resp.Text(r.Context(), w, "It Works!")
}
