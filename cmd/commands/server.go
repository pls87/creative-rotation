package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/server/http"
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
		storage := storage.New(cfg.DB)
		cr := app.New(logg, storage)

		server := http.NewServer(logg, cr, cfg.API)

		ctx, cancel := signal.NotifyContext(context.Background(),
			syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		go func() {
			<-ctx.Done()

			ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
			defer cancel()

			if err := storage.Dispose(); err != nil {
				logg.Error("failed to close storage connection: " + err.Error())
			}

			if err := server.Stop(ctx); err != nil {
				logg.Error("failed to stop http internal: " + err.Error())
			}
		}()

		logg.Info("connecting to storage...")
		var err error
		for r := retries; r > 0; r-- {
			if err = storage.Init(ctx); err == nil {
				break
			}
			logg.Errorf("failed to connect to storage: %s", err)
			logg.Info("retrying...")
			time.Sleep(retryGap)
		}

		if err != nil {
			logg.Errorf("number of retries exceeded. Shutting down: %s", err)
			cancel()
			os.Exit(1)
		}

		logg.Info("app is running...")

		if err := server.Start(ctx); err != nil {
			logg.Error("failed to start http internal: " + err.Error())
			cancel()
			os.Exit(1)
		}

		<-ctx.Done()
	},
}
