package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"

	// init postgres driver.
	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateStatsCmd)
}

var updateStatsCmd = &cobra.Command{
	Use:   "update_stats",
	Short: "Updates impressions/conversions statistics",
	Run: func(cmd *cobra.Command, args []string) {
		logg.Info("stats updater process starting...")
		defer func() {
			logg.Info("stats updater process finished...")
		}()

		var db *sqlx.DB
		var err error
		for r := retries; r > 0; r-- {
			if db, err = sqlx.Connect("postgres", cfg.DB.ConnString()); err == nil {
				break
			}
			logg.Errorf("failed to connect to storage: %s", err)
			logg.Info("retrying...")
			time.Sleep(retryGap)
		}

		if err != nil {
			logg.Errorf("number of retries exceeded. Shutting down: %s", err)
			os.Exit(1)
		}

		logg.Info("connected to storage")
		defer db.Close()

		ctx, cancel := signal.NotifyContext(context.Background(),
			syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		ticker := time.NewTicker(time.Duration(cfg.Stats.Interval) * time.Second)

		for {
			select {
			case <-ticker.C:
				err := updateStats(db)
				if err != nil {
					logg.Errorf("Couldn't update stats: %s", err)
				}
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
			logg.Errorf("Failed to rollback transaction: %s", err)
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
			logg.Errorf("Failed to rollback transaction: %s", err)
		}
		return err
	}
	return tx.Commit()
}
