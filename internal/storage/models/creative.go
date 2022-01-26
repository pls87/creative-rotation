package models

import "time"

type Creative struct {
	ID   ID     `db:"ID" json:"id"`
	Desc string `db:"description" json:"desc"`
}

type CreativeCollection struct {
	Creatives []Creative `db:"description" json:"creatives"`
}

type SlotCreative struct {
	SlotID     ID `db:"slot_id"`
	CreativeID ID `db:"creative_id"`
}

type Impression struct {
	ID         ID        `db:"ID"`
	SlotID     ID        `db:"slot_id"`
	CreativeID ID        `db:"creative_id"`
	SegmentID  ID        `db:"segment_id"`
	Time       time.Time `db:"time"`
}

type Conversion struct {
	ID         ID        `db:"ID"`
	SlotID     ID        `db:"slot_id"`
	CreativeID ID        `db:"creative_id"`
	SegmentID  ID        `db:"segment_id"`
	Time       time.Time `db:"time"`
}
