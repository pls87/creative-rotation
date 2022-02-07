package helpers

import (
	"encoding/json"
	"fmt"
)

type Entity struct {
	ID   int    `json:"id"`
	Desc string `json:"desc"`
}

func ParseOne(name string, data []byte) (ent Entity, err error) {
	m := map[string]Entity{}
	fmt.Println(string(data))
	err = json.Unmarshal(data, &m)

	return m[name], err
}

func ParseMany(name string, data []byte) (elems []Entity, err error) {
	m := map[string][]Entity{}
	fmt.Println(string(data))
	err = json.Unmarshal(data, &m)
	return m[name], err
}
