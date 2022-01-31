package models

type Segment struct {
	ID   ID     `db:"ID" json:"id"`
	Desc string `db:"description" json:"desc"`
}
