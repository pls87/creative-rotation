//go:build integration
// +build integration

package helpers

import (
	"encoding/json"
)

type StatHelper struct {
	httpHelper *HTTPHelper
}

func NewStatHelper(baseURL string) *StatHelper {
	return &StatHelper{httpHelper: NewHTTPHelper(baseURL)}
}

func (sh *StatHelper) TrackImpression(creativeID, slotID, segmentID int) (code int, body []byte, err error) {
	return sh.trackEvent("impression", creativeID, slotID, segmentID)
}

func (sh *StatHelper) TrackConversion(creativeID, slotID, segmentID int) (code int, body []byte, err error) {
	return sh.trackEvent("conversion", creativeID, slotID, segmentID)
}

func (sh *StatHelper) trackEvent(kind string, creativeID, slotID, segmentID int) (code int, body []byte, err error) {
	b, err := json.Marshal(Event{
		SlotID:     slotID,
		CreativeID: creativeID,
		SegmentID:  segmentID,
	})
	if err != nil {
		return 0, nil, err
	}
	return sh.httpHelper.Post("/"+kind, b)
}
