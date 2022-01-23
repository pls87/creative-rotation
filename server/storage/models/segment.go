package models

type Segment struct {
	ID   ID     `db:"ID"`
	Desc string `db:"description"`
}
