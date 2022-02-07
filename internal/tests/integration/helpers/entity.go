//go:build integration
// +build integration

package helpers

import "fmt"

type EntityHelper struct {
	httpHelper *HTTPHelper
}

func NewEntityHelper(baseURL string) *EntityHelper {
	return &EntityHelper{httpHelper: NewHTTPHelper(baseURL)}
}

func (ch *EntityHelper) New(t, desc string) (code int, body []byte, err error) {
	return ch.httpHelper.Post("/"+t, []byte(fmt.Sprintf(`{"desc": "%s"}`, desc)))
}

func (ch *EntityHelper) All(t string) (code int, body []byte, err error) {
	return ch.httpHelper.Get("/"+t, nil)
}
