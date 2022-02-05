package middleware

import (
	"bufio"
	"context"
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
	codeFunc func(int)
}

func WrapResponseWriter(w http.ResponseWriter, codeFunc func(int)) http.ResponseWriter {
	res := &StatusResponseWriter{ResponseWriter: w, codeFunc: codeFunc}
	_, hasHijack := w.(http.Hijacker)
	_, hasReadFrom := w.(io.ReaderFrom)

	switch {
	case hasHijack && hasReadFrom:
		return struct {
			http.ResponseWriter
			http.Hijacker
			io.ReaderFrom
		}{res, res, res}
	case hasHijack:
		return struct {
			http.ResponseWriter
			http.Hijacker
		}{res, res}
	case hasReadFrom:
		return struct {
			http.ResponseWriter
			io.ReaderFrom
		}{res, res}
	default:
		return struct {
			http.ResponseWriter
		}{res}
	}
}

func (rw *StatusResponseWriter) Header() http.Header {
	return rw.ResponseWriter.Header()
}

func (rw *StatusResponseWriter) WriteHeader(code int) {
	rw.ResponseWriter.WriteHeader(code)
	rw.codeFunc(code)
}

func (rw *StatusResponseWriter) Write(bts []byte) (n int, err error) {
	return rw.ResponseWriter.Write(bts)
}

func (rw *StatusResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return rw.ResponseWriter.(http.Hijacker).Hijack()
}

func (rw *StatusResponseWriter) ReadFrom(src io.Reader) (int64, error) {
	return rw.ResponseWriter.(io.ReaderFrom).ReadFrom(src)
}

type LoggingMiddleware struct {
	handler http.Handler
	logger  *logger.Logger
}

func (l *LoggingMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	requestID := atomic.AddInt64(&requestSeq, 1)
	ctx := context.WithValue(r.Context(), ContextKey("request_id"), requestID)
	var status *int
	wrapped := WrapResponseWriter(w, func(s int) {
		status = &s
	})
	l.handler.ServeHTTP(wrapped, r.Clone(ctx))
	l.logger.Infof(`%s [%s] %s %s %s %d "%s" request_id: %d request_time: %v`,
		r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"),
		r.Method, r.URL, r.Proto, status, r.UserAgent(), requestID, time.Since(start))
}

func NewLogger(handlerToWrap http.Handler, logger *logger.Logger) *LoggingMiddleware {
	return &LoggingMiddleware{handlerToWrap, logger}
}
