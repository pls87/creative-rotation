package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	// init postgres driver.
	_ "github.com/lib/pq"
	"github.com/pls87/creative-rotation/internal/stats"
	"github.com/pls87/creative-rotation/internal/storage"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(updateStatsCmd)
}

var updateStatsCmd = &cobra.Command{
	Use:   "update_stats",
	Short: "Updates impressions/conversions statistics",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO Clean up this long method
		logg.Info(cfg)
		logg.Info("stats updater process starting...")
		defer func() {
			logg.Info("stats updater process finished...")
		}()

		ctx, cancel := signal.NotifyContext(context.Background(),
			syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		retry := func(toRetry func() error) {
			var err error
			for r := retries; r > 0; r-- {
				if err = toRetry(); err == nil {
					break
				}
				logg.Errorf("failed to connect: %s", err)
				logg.Info("retrying...")
				time.Sleep(retryGap)
			}

			if err != nil {
				logg.Errorf("number of retries exceeded. Shutting down: %s", err)
				cancel()
				os.Exit(1)
			}
		}

		storage := storage.New(cfg.DB)
		retry(func() error {
			return storage.Init(ctx)
		})

		logg.Info("connected to storage")
		defer storage.Dispose()

		consumer := stats.NewConsumer(cfg.Queue)
		retry(func() error {
			return consumer.Init()
		})

		logg.Info("connected to rabbit")
		defer consumer.Dispose()

		logg.Info("connecting to impression queue")
		impressions, impErrors, err := consumer.Consume("StatsUpdater", stats.ImpressionQueue)
		if err != nil {
			logg.Errorf("couldn't consume impressions: %s", err)
			cancel()
			os.Exit(1)
		}

		go func() {
			for e := range impErrors {
				logg.Errorf("error while consuming impression :%s", e)
			}
		}()

		go func() {
			for e := range impressions {
				err = storage.Stats().TrackConversion(context.Background(), e.CreativeID, e.SlotID, e.SegmentID)
				if err != nil {
					logg.Errorf("couln't update impression stats by event %v: %s", e, err)
				}
			}
		}()

		logg.Info("connecting to conversion queue")
		conversions, convErrors, err := consumer.Consume("StatsUpdater", stats.ConversionQueue)
		if err != nil {
			logg.Errorf("couldn't consume impressions: %s", err)
			cancel()
			os.Exit(1)
		}

		go func() {
			for e := range convErrors {
				logg.Errorf("error while consuming conversion: %s", e)
			}
		}()

		go func() {
			for e := range conversions {
				err = storage.Stats().TrackConversion(context.Background(), e.CreativeID, e.SlotID, e.SegmentID)
				if err != nil {
					logg.Errorf("couln't update conversion stats by event %v: %s", e, err)
				}
			}
		}()

		<-ctx.Done()
	},
}
