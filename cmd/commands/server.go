package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pls87/creative-rotation/internal/server/http"

	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/storage"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(runServerCmd)
}

var runServerCmd = &cobra.Command{
	Use:   "server",
	Short: "Runs API for creative rotation app",
	Long:  `<Long version desc>`,
	Run: func(cmd *cobra.Command, args []string) {
		storage := storage.New(cfg.DB)
		cr := app.New(logg, storage, cfg)

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

		if err := storage.Init(ctx); err != nil {
			logg.Error("failed to connect to storage: " + err.Error())
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
