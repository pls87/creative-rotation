package middleware

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/pls87/creative-rotation/internal/logger"
)

var requestSeq int64

type ContextKey string

type StatusResponseWriter struct {
	http.ResponseWriter
	status int
}

func WrapResponseWriter(w http.ResponseWriter) StatusResponseWriter {
	return StatusResponseWriter{ResponseWriter: w}
}

func (rw *StatusResponseWriter) Status() int {
	return rw.status
}

func (rw *StatusResponseWriter) WriteHeader(code int) {
	if rw.Status() == 0 {
		rw.status = code
		rw.ResponseWriter.WriteHeader(code)
	}
}

func (rw *StatusResponseWriter) Write(bts []byte) (n int, err error) {
	if rw.status == 0 {
		rw.WriteHeader(http.StatusOK)
	}

	return rw.ResponseWriter.Write(bts)
}

func (rw *StatusResponseWriter) Flush() {
	if fl, ok := rw.ResponseWriter.(http.Flusher); ok {
		if rw.Status() == 0 {
			rw.WriteHeader(http.StatusOK)
		}

		fl.Flush()
	}
}

func (rw *StatusResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := rw.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, fmt.Errorf("the hijacker interface is not supported")
	}

	return hj.Hijack()
}

func (rw *StatusResponseWriter) ReadFrom(src io.Reader) (int64, error) {
	r, ok := rw.ResponseWriter.(io.ReaderFrom)
	if !ok {
		return 0, fmt.Errorf("the readerfrom interface is not supported")
	}

	return r.ReadFrom(src)
}

type LoggingMiddleware struct {
	handler http.Handler
	logger  *logger.Logger
}

func (l *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := atomic.AddInt64(&requestSeq, 1)
	ctx := context.WithValue(r.Context(), ContextKey("request_id"), requestID)
	wrapped := WrapResponseWriter(w)
	l.handler.ServeHTTP(&wrapped, r.Clone(ctx))
	l.logger.Infof(`%s [%s] %s %s %s %d "%s" request_id: %d request_time: %v`,
		r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		r.Method, r.URL, r.Proto, wrapped.Status(), r.UserAgent(), requestID, time.Since(start))
}

func NewLogger(handlerToWrap http.Handler, logger *logger.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{handlerToWrap, logger}
}
