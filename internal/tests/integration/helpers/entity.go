//go:build integration
// +build integration

package helpers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

type EntityHelper struct {
	httpHelper *HTTPHelper
}

func NewEntityHelper(baseURL string) *EntityHelper {
	return &EntityHelper{httpHelper: NewHTTPHelper(baseURL)}
}

func (ch *EntityHelper) push(t, desc string) (code int, body []byte, err error) {
	return ch.httpHelper.Post("/"+t, []byte(fmt.Sprintf(`{"desc": "%s"}`, desc)))
}

func (ch *EntityHelper) getAll(t string) (code int, body []byte, err error) {
	return ch.httpHelper.Get("/"+t, nil)
}

func (ch *EntityHelper) New(t *testing.T, kind, desc string) (entity Entity) {
	code, resp, err := ch.push(kind, desc)

	require.NoErrorf(t, err, "no error expected but was: %s", err)
	require.Equal(t, http.StatusOK, code)
	entity, err = ParseOne(kind, resp)
	require.NoErrorf(t, err, "no error expected but was: %s", err)

	return entity
}

func (ch *EntityHelper) All(t *testing.T, kind string) (entities []Entity) {
	code, resp, err := ch.getAll(kind)
	require.NoErrorf(t, err, "no error expected but was: %s", err)
	require.Equal(t, http.StatusOK, code)
	entities, err = ParseMany(kind+"s", resp)
	require.NoErrorf(t, err, "no error expected but was: %s", err)

	return entities
}