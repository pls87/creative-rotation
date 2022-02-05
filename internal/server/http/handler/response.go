package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pls87/creative-rotation/internal/logger"
)

type response struct {
	logger *logger.Logger
}

func (eh *response) httpError(ctx context.Context, w http.ResponseWriter, status int, msg string, err error) {
	eh.logger.WithContext(ctx).Error(fmt.Errorf("%s: %w", msg, err))
	http.Error(w, msg, status)
}

func (eh *response) badRequest(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	eh.httpError(ctx, w, http.StatusBadRequest, msg, err)
}

func (eh *response) internalServerError(ctx context.Context, w http.ResponseWriter, msg string, err error) {
	eh.httpError(ctx, w, http.StatusInternalServerError, msg, err)
}

func (eh *response) json(ctx context.Context, w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		eh.internalServerError(ctx, w, "failed to encode response body", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(data)
	if err != nil {
		eh.internalServerError(ctx, w, "failed to write response body", err)
	}
}

func (eh *response) text(ctx context.Context, w http.ResponseWriter, msg string) {
	if _, err := w.Write([]byte(msg)); err != nil {
		eh.internalServerError(ctx, w, "failed to write response body", err)
	}
}
