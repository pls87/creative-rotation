package models

type Stats struct {
	SlotID      ID     `db:"slot_id"`
	CreativeID  ID     `db:"creative_id"`
	SegmentID   ID     `db:"segment_id"`
	Impressions uint64 `db:"impressions"`
	Conversions uint64 `db:"conversions"`
}
