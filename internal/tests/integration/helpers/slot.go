//go:build integration
// +build integration

package helpers

import "fmt"

type SlotHelper struct {
	httpHelper *HTTPHelper
}

func NewSlotHelper(baseURL string) *SlotHelper {
	return &SlotHelper{httpHelper: NewHTTPHelper(baseURL)}
}

func (sh *SlotHelper) Creatives(slotID int) (code int, body []byte, err error) {
	url := fmt.Sprintf("/slot/%d/creative", slotID)
	return sh.httpHelper.Get(url, nil)
}
