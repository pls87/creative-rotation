package models

type Slot struct {
	ID   ID     `db:"ID"`
	Desc string `db:"description"`
}
