package models

import "time"

type Creative struct {
	ID   ID     `db:"ID" json:"id"`
	Desc string `db:"description" json:"desc"`
}

type SlotCreative struct {
	SlotID     ID `db:"slot_id" json:"slot_id"`
	CreativeID ID `db:"creative_id" json:"creative_id"`
}

type Impression struct {
	ID         ID        `db:"ID" json:"id"`
	SlotID     ID        `db:"slot_id" json:"slot_id"`
	CreativeID ID        `db:"creative_id" json:"creative_id"`
	SegmentID  ID        `db:"segment_id" json:"segment_id"`
	Time       time.Time `db:"time" json:"time"`
}

type Conversion struct {
	ID         ID        `db:"ID"`
	SlotID     ID        `db:"slot_id" json:"slot_id"`
	CreativeID ID        `db:"creative_id" json:"creative_id"`
	SegmentID  ID        `db:"segment_id" json:"segment_id"`
	Time       time.Time `db:"time" json:"time"`
}
