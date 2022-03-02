//go:build integration
// +build integration

package helpers

import (
	"encoding/json"
)

type Entity struct {
	ID   int    `json:"id"`
	Desc string `json:"desc"`
}

type SlotCreative struct {
	SlotID       int    `json:"slot_id"`
	SlotDesc     string `json:"slot_desc"`
	CreativeID   int    `json:"creative_id"`
	CreativeDesc string `json:"creative_desc"`
}

type Event struct {
	SlotID     int `json:"slot_id"`
	CreativeID int `json:"creative_id"`
	SegmentID  int `json:"segment_id"`
}

type Stats struct {
	SlotID      int `db:"slot_id"`
	CreativeID  int `db:"creative_id"`
	SegmentID   int `db:"segment_id"`
	Impressions int `db:"impressions"`
	Conversions int `db:"conversions"`
}

func ParseOne(data []byte) (ent Entity, err error) {
	err = json.Unmarshal(data, &ent)

	return ent, err
}

func ParseMany(name string, data []byte) (elems []Entity, err error) {
	m := map[string][]Entity{}
	err = json.Unmarshal(data, &m)
	return m[name], err
}
