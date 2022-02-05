package http

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/config"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/server/http/handler"
	"github.com/pls87/creative-rotation/internal/server/http/middleware"
)

type Server struct {
	httpServer *http.Server
	service    *handler.Service
	cfg        config.APIConf
	logger     *logger.Logger
}

func NewServer(logger *logger.Logger, app app.Application, cfg config.APIConf) *Server {
	return &Server{
		logger:  logger,
		cfg:     cfg,
		service: handler.NewService(app, logger),
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := mux.NewRouter()

	mux.HandleFunc("/noop", s.service.Noop).Methods("GET")
	mux.HandleFunc("/creative", s.service.Creatives().All).Methods("GET")
	mux.HandleFunc("/creative", s.service.Creatives().New).Methods("POST")
	mux.HandleFunc("/creative/{creative_id:[0-9]+}/slot", s.service.Creatives().AddToSlot).
		Methods("POST")
	mux.HandleFunc("/creative/{creative_id:[0-9]+}/slot/{slot_id:[0-9]+}", s.service.Creatives().RemoveFromSlot).
		Methods("DELETE")
	mux.HandleFunc("/conversion", s.service.Creatives().TrackConversion).Methods("POST")
	mux.HandleFunc("/impression", s.service.Creatives().TrackImpression).Methods("POST")
	mux.HandleFunc("/creative/next", s.service.Creatives().Next).Methods("GET")
	mux.HandleFunc("/slot", s.service.Slots().All).Methods("GET")
	mux.HandleFunc("/slot", s.service.Slots().New).Methods("POST")
	mux.HandleFunc("/segment", s.service.Segments().All).Methods("GET")
	mux.HandleFunc("/segment", s.service.Segments().New).Methods("POST")

	s.httpServer = &http.Server{
		Addr:    net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port)),
		Handler: middleware.NewLogger(mux, s.logger),
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
