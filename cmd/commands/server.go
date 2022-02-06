package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/server/http"
	"github.com/pls87/creative-rotation/internal/stats"
	"github.com/pls87/creative-rotation/internal/storage"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runServerCmd)
}

var runServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs REST API for creative rotation app",
	Run: func(cmd *cobra.Command, args []string) {
		// TODO Clean up this long method
		logg.Info(cfg)
		storage := storage.New(cfg.DB)
		stats := stats.NewProducer(cfg.Queue)
		cr := app.New(logg, storage, stats)

		server := http.NewServer(logg, cr, cfg.API)

		ctx, cancel := signal.NotifyContext(context.Background(),
			syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		shutDown := func() {
			<-ctx.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			if err := storage.Dispose(); err != nil {
				logg.Error("failed to close storage connection: " + err.Error())
			}

			if err := stats.Dispose(); err != nil {
				logg.Error("failed to close rabbit connection: " + err.Error())
			}

			if err := server.Stop(ctx); err != nil {
				logg.Error("failed to stop http internal: " + err.Error())
			}
		}

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

		logg.Info("connecting to storage...")
		retry(func() error {
			return storage.Init(ctx)
		})

		logg.Info("connecting to rabbit...")
		retry(func() error {
			return stats.Init()
		})

		logg.Info("app is running...")

		go shutDown()

		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http server: " + err.Error())
			cancel()
			os.Exit(1)
		}
	},
}
