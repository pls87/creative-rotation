package sql

import (
	"fmt"

	"github.com/pls87/creative-rotation/internal/storage/models"
)

func ALLQuery(entity string) (query string) {
	return fmt.Sprintf(`SELECT * FROM "%s"`, entity)
}

func CREATEQuery(entity string) (query string) {
	return fmt.Sprintf(`INSERT INTO "%s" (description) VALUES ($1) RETURNING "ID"`, entity)
}

func DELETEQuery(entity string) (query string) {
	return fmt.Sprintf(`DELETE FROM "%s" WHERE "ID"=$1`, entity)
}

func ALLError(entity string, err error) error {
	return fmt.Errorf("couldn't get %ss from database: %w", entity, err)
}

func CREATEError(entity string, err error) error {
	return fmt.Errorf("couldn't create %s in database: %w", entity, err)
}

func DELETEError(entity string, id models.ID, err error) error {
	return fmt.Errorf("couldn't delete %s with id=%d: %w", entity, id, err)
}
