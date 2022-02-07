package errors

import (
	"fmt"

	"github.com/pls87/creative-rotation/internal/storage/models"
)

type SlotCreativeErrors struct{}

func (c *SlotCreativeErrors) GetFor(primary, secondary string, id models.ID, err error) error {
	return fmt.Errorf("couldn't get %ss for %s '%d' from database: %w", secondary, primary, id, err)
}

func (c *SlotCreativeErrors) Create(creativeID, slotID models.ID, err error) error {
	return fmt.Errorf("couldn't add creative id=%d to slot_id=%d: %w", creativeID, slotID, err)
}

func (c *SlotCreativeErrors) Delete(creativeID, slotID models.ID, err error) error {
	return fmt.Errorf("couldn't delete creative id=%d from slot_id=%d: %w", creativeID, slotID, err)
}

func (c *SlotCreativeErrors) Exists(creativeID, slotID models.ID, err error) error {
	return fmt.Errorf("couldn't check if creative id=%d added into slot_id=%d: %w", creativeID, slotID, err)
}

var SC = SlotCreativeErrors{}
