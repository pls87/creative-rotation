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

func Parse(data []byte, target interface{}) error {
	return json.Unmarshal(data, &target)
}

func ParseOne(data []byte) (ent Entity, err error) {
	err = Parse(data, &ent)

	return ent, err
}

func ParseMany(name string, data []byte) (elems []Entity, err error) {
	m := map[string][]Entity{}
	err = Parse(data, &m)
	return m[name], err
}
