package queries

import "fmt"

const SlotCreativeRelation = "slot_creative"

type SlotCreativeQueries struct{}

func (l *SlotCreativeQueries) GetFor(primary, secondary string) string {
	return fmt.Sprintf(`SELECT s.* FROM "%s" sc
		INNER JOIN "%s" s ON sc.%s_id=s."ID"
		WHERE sc.%s_id = $1`, SlotCreativeRelation, secondary, secondary, primary)
}

func (l *SlotCreativeQueries) Create() string {
	return fmt.Sprintf(`INSERT INTO "%s" (creative_id, slot_id) VALUES ($1, $2)`, SlotCreativeRelation)
}

func (l *SlotCreativeQueries) Exists() string {
	return fmt.Sprintf(`SELECT * FROM "%s" WHERE creative_id = $1 AND slot_id = $2`, SlotCreativeRelation)
}

func (l *SlotCreativeQueries) Delete() string {
	return fmt.Sprintf(`DELETE FROM "%s" WHERE creative_id = $1 AND slot_id=$2`, SlotCreativeRelation)
}

func (l *SlotCreativeQueries) All() string {
	return fmt.Sprintf(`SELECT sc.slot_id, sc.creative_id, s.description as slot_desc, 
	cr.description as creative_desc 
		FROM "%s" sc 
		INNER JOIN "%s" s ON sc.slot_id=s."ID"
		INNER JOIN "%s" cr ON sc.creative_id=cr."ID"`, SlotCreativeRelation, SlotRelation, CreativeRelation)
}

var SC = SlotCreativeQueries{}
