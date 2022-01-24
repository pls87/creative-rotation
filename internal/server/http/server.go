package http

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/pls87/creative-rotation/internal/app"

	"github.com/pls87/creative-rotation/internal/config"
	"github.com/pls87/creative-rotation/internal/logger"
)

type Server struct {
	app         app.Application
	httpServer  *http.Server
	httpService *Service
	cfg         config.APIConf
	logger      *logger.Logger
}

func NewServer(logger *logger.Logger, app app.Application, cfg config.APIConf) *Server {
	return &Server{
		logger:      logger,
		app:         app,
		cfg:         cfg,
		httpService: NewService(logger),
	}
}

func (s *Server) Start(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/noop", s.httpService.Noop)

	s.httpServer = &http.Server{
		Addr:    net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port)),
		Handler: NewLogger(mux, s.logger),
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
