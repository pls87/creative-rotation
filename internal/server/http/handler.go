package http

import (
	"net/http"

	"github.com/pls87/creative-rotation/internal/logger"
)

type Service struct {
	logger *logger.Logger
}

type wrappedResponseWriter struct {
	http.ResponseWriter
	status int
}

func wrapResponseWriter(w http.ResponseWriter) wrappedResponseWriter {
	return wrappedResponseWriter{ResponseWriter: w}
}

func (rw *wrappedResponseWriter) Status() int {
	return rw.status
}

func (rw *wrappedResponseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

func NewService(logger *logger.Logger) *Service {
	return &Service{
		logger: logger,
	}
}

func (s *Service) Noop(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("It works!"))
}
