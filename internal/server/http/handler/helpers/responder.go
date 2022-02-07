package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pls87/creative-rotation/internal/logger"
)

type Responder struct {
	logger *logger.Logger
}

func NewResponder(l *logger.Logger) *Responder {
	return &Responder{logger: l}
}

func (eh *Responder) httpError(ctx context.Context, w http.ResponseWriter, status int, msg string, err error) {
	eh.logger.WithContext(ctx).Error(fmt.Errorf("%s: %w", msg, err))
	http.Error(w, msg, status)
}

func (eh *Responder) BadRequest(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	eh.httpError(ctx, w, http.StatusBadRequest, msg, err)
}

func (eh *Responder) InternalServerError(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	eh.httpError(ctx, w, http.StatusInternalServerError, msg, err)
}

func (eh *Responder) JSON(ctx context.Context, w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		eh.InternalServerError(ctx, w, "failed to encode response body", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		eh.InternalServerError(ctx, w, "failed to write response body", err)
	}
}

func (eh *Responder) Text(ctx context.Context, w http.ResponseWriter, msg string) {
	if _, err := w.Write([]byte(msg)); err != nil {
		eh.InternalServerError(ctx, w, "failed to write response body", err)
	}
}
