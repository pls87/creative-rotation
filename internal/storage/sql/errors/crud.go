package errors

import (
	"fmt"

	"github.com/pls87/creative-rotation/internal/storage/models"
)

type CrudErrors struct{}

func (c *CrudErrors) All(entity string, err error) error {
	return fmt.Errorf("couldn't get %ss from database: %w", entity, err)
}

func (c *CrudErrors) Create(entity string, err error) error {
	return fmt.Errorf("couldn't create %s in database: %w", entity, err)
}

func (c *CrudErrors) Delete(entity string, id models.ID, err error) error {
	return fmt.Errorf("couldn't delete %s with id=%d: %w", entity, id, err)
}

var Crud = CrudErrors{}
