package queries

import (
	"fmt"

	"github.com/pls87/creative-rotation/internal/storage/models"
)

const (
	StatsRelation   = "stats"
	ImpressionField = "impressions"
	ConversionField = "conversions"
)

type StatsQueries struct{}

func (l *StatsQueries) For() string {
	return fmt.Sprintf(`SELECT sc.slot_id as slot_id, 
		sc.creative_id as creative_id, seg."ID" as segment_id,
		CASE WHEN st.impressions is NULL THEN 0 ELSE st.impressions END as impressions,
		CASE WHEN st.conversions is NULL THEN 0 ELSE st.conversions END as conversions
		FROM "%s" sc CROSS JOIN "%s" seg LEFT JOIN "%s" st ON sc.creative_id = st.creative_id 
		AND sc.slot_id = st.slot_id AND seg."ID" = st.segment_id
		WHERE sc.slot_id = $1 and seg."ID" = $2`, SlotCreativeRelation, SegmentRelation, StatsRelation)
}

func (l *StatsQueries) Track(field string, cid, slid, segid models.ID) string {
	return l.updateStats(field, cid, slid, segid)
}

// Standard sql parameters substitution didn't work for such multi-level query.
func (l *StatsQueries) updateStats(field string, cid, slid, segid models.ID) string {
	return fmt.Sprintf(`INSERT INTO "%s" (%s, creative_id, slot_id, segment_id) 
		VALUES (1, %d, %d, %d)
	ON CONFLICT (creative_id, slot_id, segment_id) 
	DO UPDATE SET %s = (SELECT %s FROM "%s" WHERE creative_id=%d AND slot_id=%d AND segment_id=%d) + 1`,
		StatsRelation, field, cid, slid, segid, field, field, StatsRelation, cid, slid, segid)
}

var Stats = StatsQueries{}
