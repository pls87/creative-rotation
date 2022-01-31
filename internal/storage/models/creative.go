package models

import "time"

type Creative struct {
	ID   ID     `db:"ID" json:"id"`
	Desc string `db:"description" json:"desc"`
}

type SlotCreative struct {
	SlotID     ID `db:"slot_id" json:"slot_id"`         //nolint: tagliatelle
	CreativeID ID `db:"creative_id" json:"creative_id"` //nolint: tagliatelle
}

type Impression struct {
	ID         ID        `db:"ID" json:"id"`
	SlotID     ID        `db:"slot_id" json:"slot_id"`         //nolint: tagliatelle
	CreativeID ID        `db:"creative_id" json:"creative_id"` //nolint: tagliatelle
	SegmentID  ID        `db:"segment_id" json:"segment_id"`   //nolint: tagliatelle
	Time       time.Time `db:"time" json:"time"`
}

type Conversion struct {
	ID         ID        `db:"ID"`
	SlotID     ID        `db:"slot_id" json:"slot_id"`         //nolint: tagliatelle
	CreativeID ID        `db:"creative_id" json:"creative_id"` //nolint: tagliatelle
	SegmentID  ID        `db:"segment_id" json:"segment_id"`   //nolint: tagliatelle
	Time       time.Time `db:"time" json:"time"`
}
