package commands

import (
	"context"
	"errors"
	"net/http"
	"os"
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
	ctx     context.Context
	cancel  context.CancelFunc
}

func (sc *ServerCMD) onFail() {
	sc.cancel()
	os.Exit(1)
}

func (sc *ServerCMD) shutDown() {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := sc.storage.Dispose(); err != nil {
		sc.logg.Error("failed to close storage connection: " + err.Error())
	}

	if err := sc.stats.Dispose(); err != nil {
		sc.logg.Error("failed to close rabbit connection: " + err.Error())
	}

	if err := sc.server.Stop(ctx); err != nil {
		sc.logg.Error("failed to stop http server: " + err.Error())
	}
}

func (sc *ServerCMD) Run() {
	sc.storage = storage.New(sc.cfg.DB)
	sc.stats = stats.NewProducer(sc.cfg.Queue)
	sc.cr = app.New(sc.logg, sc.storage, sc.stats)
	sc.server = httpsrv.NewServer(sc.logg, sc.cr, sc.cfg.API)

	sc.ctx, sc.cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer sc.cancel()

	sc.logg.Info("connecting to storage...")
	sc.Retry(sc.ctx, func() error {
		return sc.storage.Init(sc.ctx)
	}, sc.onFail)

	sc.logg.Info("connecting to rabbit...")
	sc.Retry(sc.ctx, func() error {
		return sc.stats.Init()
	}, sc.onFail)

	go func() {
		if err := sc.server.Start(sc.ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
			sc.logg.Error("server was unexpectedly stopped: " + err.Error())
			sc.onFail()
		}
	}()

	<-sc.ctx.Done()

	sc.shutDown()
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
