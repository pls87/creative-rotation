package http

import (
	"context"
	"net"
	"net/http"
	"strconv"

	"github.com/pls87/creative-rotation/internal/server/http/handler"
	"github.com/pls87/creative-rotation/internal/server/http/middleware"

	mux2 "github.com/gorilla/mux"

	"github.com/pls87/creative-rotation/internal/app"

	"github.com/pls87/creative-rotation/internal/config"
	"github.com/pls87/creative-rotation/internal/logger"
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
	mux := mux2.NewRouter()

	mux.HandleFunc("/noop", s.service.Noop)
	mux.HandleFunc("/creatives", s.service.Creatives().All)

	s.httpServer = &http.Server{
		Addr:    net.JoinHostPort(s.cfg.Host, strconv.Itoa(s.cfg.Port)),
		Handler: middleware.NewLogger(mux, s.logger),
	}

	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
