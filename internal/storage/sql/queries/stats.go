package queries

import "fmt"

const (
	StatsRelation   = "stats"
	ImpressionField = "impressions"
	ConversionField = "conversions"
)

type StatsQueries struct{}

func (l *StatsQueries) For() string {
	return fmt.Sprintf(`SELECT * FROM "%s" WHERE slot_id=$1 AND segment_id=$2`, StatsRelation)
}

func (l *StatsQueries) Track(field string) string {
	return l.updateStats(field)
}

func (l *StatsQueries) updateStats(field string) string {
	return fmt.Sprintf(`INSERT INTO "%s" (%s, creative_id, slot_id, segment_id) 
		SELECT * FROM  (VALUES (1, $1, $2, $3)) AS t(%s, creative_id, slot_id, segment_id)
		WHERE EXISTS (
			SELECT FROM "slot_creative" sc
			WHERE  sc.slot_id = t.slot_id
				AND    sc.creative_id = t.creative_id
		)
		ON CONFLICT (creative_id, slot_id, segment_id) DO UPDATE SET %s = EXCLUDED.%s + 1`,
		StatsRelation, field, field, field, field)
}

var Stats = StatsQueries{}
