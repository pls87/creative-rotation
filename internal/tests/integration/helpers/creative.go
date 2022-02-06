//go:build integration
// +build integration

package helpers

import (
	"fmt"
	"strconv"
)

type CreativeHelper struct {
	httpHelper *HTTPHelper
}

func NewCreativeHelper(baseURL string) *CreativeHelper {
	return &CreativeHelper{httpHelper: NewHTTPHelper(baseURL)}
}

func (ch *CreativeHelper) New(desc string) (code int, body []byte, err error) {
	return ch.httpHelper.Post("/creative", []byte(fmt.Sprintf(`{"desc": "%s"}`, desc)))
}

func (ch *CreativeHelper) All(desc string) (code int, body []byte, err error) {
	return ch.httpHelper.Get("/creative", nil)
}

func (ch *CreativeHelper) AddToSlot(creativeID, slotID int) (code int, body []byte, err error) {
	url := fmt.Sprintf("/creative/%d/slot", creativeID)
	slotBody := []byte(fmt.Sprintf(`{"id:"%d"}`, slotID))
	return ch.httpHelper.Post(url, slotBody)
}

func (ch *CreativeHelper) RemoveFromSlot(creativeID, slotID int) (code int, body []byte, err error) {
	url := fmt.Sprintf("/creative/%d/slot/%d", creativeID, slotID)
	return ch.httpHelper.Delete(url, nil)
}

func (ch *CreativeHelper) Next(slotID, segmentID int) (code int, body []byte, err error) {
	return ch.httpHelper.Get("/creative/next", map[string]string{
		"slot_id":    strconv.Itoa(slotID),
		"segment_id": strconv.Itoa(segmentID),
	})
}
