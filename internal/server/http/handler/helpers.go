package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/pls87/creative-rotation/internal/storage/models"
)

type helpers struct {
	resp *response
}

func (h *helpers) handleIDQuery(ctx context.Context, param string, w http.ResponseWriter,
	r *http.Request) (models.ID, bool) {
	IDs := r.URL.Query()[param]
	if len(IDs) == 0 {
		h.resp.badRequest(ctx, w, fmt.Sprintf("%s isn't specified", param), nil)
		return 0, false
	}
	if len(IDs) != 1 {
		h.resp.badRequest(ctx, w, fmt.Sprintf("more than one %s were passed", param), nil)
		return 0, false
	}

	ID, err := strconv.Atoi(IDs[0])
	if err != nil || ID <= 0 {
		h.resp.badRequest(ctx, w, fmt.Sprintf("malformed %s", param), err)
		return 0, false
	}

	return models.ID(ID), true
}

func (h *helpers) handleURLParamID(w http.ResponseWriter, r *http.Request, name string) (id models.ID, ok bool) {
	vars := mux.Vars(r)

	tempID, e := strconv.Atoi(vars[name])
	if e != nil || tempID <= 0 {
		h.resp.badRequest(r.Context(), w, fmt.Sprintf("malformed %s", name), e)
		return 0, false
	}

	return models.ID(tempID), true
}
