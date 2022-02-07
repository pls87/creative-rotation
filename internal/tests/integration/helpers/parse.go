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

func ParseOne(data []byte) (ent Entity, err error) {
	err = json.Unmarshal(data, &ent)

	return ent, err
}

func ParseMany(name string, data []byte) (elems []Entity, err error) {
	m := map[string][]Entity{}
	err = json.Unmarshal(data, &m)
	return m[name], err
}
