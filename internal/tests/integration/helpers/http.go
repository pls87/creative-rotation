//go:build integration
// +build integration

package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type HTTPHelper struct {
	client  http.Client
	baseURL string
}

func NewHTTPHelper(baseURL string) *HTTPHelper {
	return &HTTPHelper{
		client: http.Client{
			Timeout: 2 * time.Second,
		},
		baseURL: baseURL,
	}
}

func (h *HTTPHelper) Post(url string, body []byte) (code int, respBody []byte, err error) {
	return h.sendRequest(http.MethodPost, url, nil, body)
}

func (h *HTTPHelper) Delete(url string, body []byte) (code int, respBody []byte, err error) {
	return h.sendRequest(http.MethodDelete, url, nil, body)
}

func (h *HTTPHelper) Get(url string, query map[string]string) (code int, respBody []byte, err error) {
	return h.sendRequest(http.MethodGet, url, query, nil)
}

func (h *HTTPHelper) buildURL(url string, query map[string]string) string {
	urlBuilder := strings.Builder{}
	urlBuilder.WriteString(h.baseURL)
	urlBuilder.WriteString(url)
	if len(query) == 0 {
		return urlBuilder.String()
	}
	urlBuilder.WriteString("?")
	for k, v := range query {
		urlBuilder.WriteString(fmt.Sprintf("%s=%s&", k, v))
	}

	return strings.TrimRight(urlBuilder.String(), "&")
}

func (h *HTTPHelper) sendRequest(method, url string, query map[string]string,
	body []byte) (code int, respBody []byte, err error) {
	req, err := http.NewRequest(method, h.buildURL(url, query), bytes.NewReader(body)) // nolint: noctx
	if err != nil {
		return 0, nil, err
	}
	resp, err := h.client.Do(req)
	if err != nil {
		return 0, nil, err
	}

	return h.handleResponse(resp)
}

func (h *HTTPHelper) handleResponse(resp *http.Response) (code int, respBody []byte, err error) {
	respBody = make([]byte, 4096)
	_, err = resp.Body.Read(respBody)
	if err != nil && !errors.Is(err, io.EOF) {
		return 0, nil, err
	}

	return resp.StatusCode, respBody, nil
}
