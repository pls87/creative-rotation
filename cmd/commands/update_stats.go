package commands

import (
	"context"
	"os/signal"
	"syscall"
	"time"
	// init postgres driver.
	_ "github.com/lib/pq"

	"github.com/jmoiron/sqlx"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateStatsCmd)
}

var updateStatsCmd = &cobra.Command{
	Use:   "update_stats",
	Short: "Updates impressions/conversions statistics",
	Long:  `<Long version desc>`,
	Run: func(cmd *cobra.Command, args []string) {
		logg.Infof("[%s] stats updater process starting...", time.Now())
		defer func() {
			logg.Infof("[%s] stats updater process finished...", time.Now())
		}()

		db, err := sqlx.Connect("postgres", cfg.DB.ConnString())
		if err != nil {
			logg.Errorf("[%s] Couldn't connect to database to update stats: %s", time.Now(), err)
		}
		defer db.Close()

		ctx, cancel := signal.NotifyContext(context.Background(),
			syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		ticker := time.NewTicker(time.Duration(cfg.Stats.Interval) * time.Second)

		for {
			select {
			case <-ticker.C:
				logg.Infof("[%s]Update Start...", time.Now())
				err := updateStats(db)
				if err != nil {
					logg.Errorf("[%s] ERROR: Couldn't uodate stats: %s", time.Now(), err)
				}
				logg.Infof("[%s]Update Finished!", time.Now())
			case <-ctx.Done():
				return
			}
		}
	},
}

func updateStats(db *sqlx.DB) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`TRUNCATE TABLE "stats"`)
	if err != nil {
		e := tx.Rollback()
		if e != nil {
			logg.Errorf("[%s] ERROR: Failed to rollback transaction: %s", time.Now(), err)
		}
		return err
	}

	query := `
	INSERT INTO stats (slot_id, creative_id, segment_id, impressions, conversions)
	SELECT sc.slot_id as slot_id, sc.creative_id as creative_id, s."ID" as segment_id, 
	       COUNT(DISTINCT i."ID") as impressions,COUNT(DISTINCT c."ID") as conversions FROM "slot_creative" sc
	LEFT JOIN "segment" s ON TRUE
	LEFT JOIN "impression" i ON sc.slot_id=i.slot_id AND sc.creative_id=i.creative_id AND i.segment_id=s."ID"
	LEFT JOIN "conversion" c ON sc.slot_id=c.slot_id AND sc.creative_id=c.creative_id AND c.segment_id=s."ID"
	GROUP BY (sc.slot_id, sc.creative_id, s."ID")`

	_, err = tx.Exec(query)

	if err != nil {
		e := tx.Rollback()
		if e != nil {
			logg.Errorf("[%s]ERROR: Failed to rollback transaction: %s", time.Now(), err)
		}
		return err
	}
	return tx.Commit()
}
