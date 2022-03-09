package commands

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/pls87/creative-rotation/internal/app"
	"github.com/pls87/creative-rotation/internal/config"
	"github.com/pls87/creative-rotation/internal/logger"
	httpsrv "github.com/pls87/creative-rotation/internal/server/http"
	"github.com/pls87/creative-rotation/internal/stats"
	"github.com/pls87/creative-rotation/internal/storage"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/spf13/cobra"
)

type ServerCMD struct {
	*RootCMD
	storage basic.Storage
	cr      app.Application
	stats   stats.Producer
	server  *httpsrv.Server
}

func (sc *ServerCMD) shutDown() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := sc.storage.Dispose(); err != nil {
		sc.logg.Errorf("failed to close storage connection: %v", err)
	}

	if err := sc.stats.Dispose(); err != nil {
		sc.logg.Errorf("failed to close rabbit connection: %v", err)
	}

	if err := sc.server.Stop(ctx); err != nil {
		sc.logg.Errorf("failed to stop http server: %v", err)
	}
}

func (sc *ServerCMD) Run() {
	defer sc.shutDown()

	sc.storage = storage.New(sc.cfg.DB)
	sc.stats = stats.NewProducer(sc.cfg.Queue)
	sc.cr = app.New(sc.logg, sc.storage, sc.stats)
	sc.server = httpsrv.NewServer(sc.logg, sc.cr, sc.cfg.API)

	ctx, _ := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	sc.logg.Info("connecting to storage...")
	if err := sc.Retry(ctx, func() error {
		return sc.storage.Init(ctx)
	}); err != nil {
		sc.logg.Errorf("couldn't connect to storage: %v", err)
		return
	}

	sc.logg.Info("connecting to rabbit...")
	if err := sc.Retry(ctx, func() error {
		return sc.stats.Init()
	}); err != nil {
		sc.logg.Errorf("couldn't connect to rabbit: %v", err)
		return
	}

	go func() {
		if err := sc.server.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			sc.logg.Errorf("server was unexpectedly stopped: %v", err)
			return
		}
	}()

	<-ctx.Done()
}

func (sc *ServerCMD) Init() {
	sc.cfg = config.New(sc.cfgFile)
	sc.logg = logger.New(sc.cfg.Log)
}

func newServerCmd() *ServerCMD {
	cmd := &ServerCMD{
		RootCMD: &RootCMD{},
	}
	cmd.Command = &cobra.Command{
		Use:   "server",
		Short: "Runs http api for creative rotations app",
		Run: func(c *cobra.Command, args []string) {
			cmd.Run()
		},
	}

	return cmd
}

func init() {
	cmd := newServerCmd()
	cobra.OnInitialize(cmd.Init)
	cmd.PersistentFlags().StringVar(&cmd.cfgFile, "config", "", "config file")
	rc.AddCommand(cmd.Command)
}
