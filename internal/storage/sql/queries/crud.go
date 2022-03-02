package queries

import "fmt"

const (
	SlotRelation     = "slot"
	CreativeRelation = "creative"
	SegmentRelation  = "segment"
)

type CrudQueries struct{}

func (c *CrudQueries) All(entity string) (query string) {
	return fmt.Sprintf(`SELECT * FROM "%s"`, entity)
}

func (c *CrudQueries) Create(entity string) (query string) {
	return fmt.Sprintf(`INSERT INTO "%s" ("description") VALUES ($1) RETURNING "ID"`, entity)
}

func (c *CrudQueries) Delete(entity string) (query string) {
	return fmt.Sprintf(`DELETE FROM "%s" WHERE "ID"=$1`, entity)
}

var Crud = CrudQueries{}
