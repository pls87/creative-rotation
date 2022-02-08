package commands

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	// init postgres driver.
	_ "github.com/lib/pq"
	"github.com/pls87/creative-rotation/internal/config"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/pls87/creative-rotation/internal/stats"
	"github.com/pls87/creative-rotation/internal/storage"
	"github.com/pls87/creative-rotation/internal/storage/basic"
	"github.com/spf13/cobra"
)

type StatsCMD struct {
	*RootCMD
	storage basic.Storage
	stats   stats.Consumer
	ctx     context.Context
	cancel  context.CancelFunc
}

func (sc *StatsCMD) onFail() {
	sc.cancel()
	os.Exit(1)
}

func (sc *StatsCMD) consumeImpressions() (chan stats.Event, chan error) {
	impressions, impErrors, err := sc.stats.Consume("StatsUpdater", stats.ImpressionQueue)
	if err != nil {
		sc.logg.Errorf("couldn't consume impressions: %s", err)
		sc.onFail()
	}

	return impressions, impErrors
}

func (sc *StatsCMD) consumeConversions() (chan stats.Event, chan error) {
	conversions, impErrors, err := sc.stats.Consume("StatsUpdater", stats.ConversionQueue)
	if err != nil {
		sc.logg.Errorf("couldn't consume conversions: %s", err)
		sc.onFail()
	}

	return conversions, impErrors
}

func (sc *StatsCMD) waitForMessages(i, c chan stats.Event, ie, ce chan error) {
	var e error
	var ev stats.Event
	for {
		select {
		case e = <-ie:
			sc.logg.Errorf("error while consuming impression: %s", e)
		case e = <-ce:
			sc.logg.Errorf("error while consuming conversion: %s", e)
		case ev = <-i:
			e = sc.storage.Stats().TrackImpression(context.Background(), ev.CreativeID, ev.SlotID, ev.SegmentID)
			if e != nil {
				sc.logg.Errorf("couln't update impression stats by event %v: %s", ev, e)
			}
		case ev = <-c:
			e = sc.storage.Stats().TrackConversion(context.Background(), ev.CreativeID, ev.SlotID, ev.SegmentID)
			if e != nil {
				sc.logg.Errorf("couln't update conversion stats by event %v: %s", ev, e)
			}
		case <-sc.ctx.Done():
			return
		}
	}
}

func (sc *StatsCMD) shutDown() {
	if err := sc.storage.Dispose(); err != nil {
		sc.logg.Error("failed to close storage connection: " + err.Error())
	}

	if err := sc.stats.Dispose(); err != nil {
		sc.logg.Error("failed to close rabbit connection: " + err.Error())
	}
}

func (sc *StatsCMD) Run() {
	sc.logg.Info("stats updater process starting...")
	defer func() {
		sc.logg.Info("stats updater process finished...")
	}()

	sc.storage = storage.New(sc.cfg.DB)
	sc.stats = stats.NewConsumer(sc.cfg.Queue)

	sc.ctx, sc.cancel = signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer sc.cancel()

	sc.logg.Info("connecting to storage...")
	sc.Retry(func() error {
		return sc.storage.Init(sc.ctx)
	}, sc.onFail)

	defer sc.storage.Dispose()

	sc.logg.Info("connecting to rabbit...")
	sc.Retry(func() error {
		return sc.stats.Init()
	}, sc.onFail)

	defer sc.stats.Dispose()

	// could get all events via one queue but would like to play with rabbit routing
	i, ie := sc.consumeImpressions()
	c, ce := sc.consumeConversions()

	sc.waitForMessages(i, c, ie, ce)

	sc.shutDown()
}

func (sc *StatsCMD) Init() {
	sc.cfg = config.New(sc.cfgFile)
	sc.logg = logger.New(sc.cfg.Log)
}

func newStatsCmd() *StatsCMD {
	cmd := StatsCMD{
		RootCMD: &RootCMD{},
	}

	cmd.Command = &cobra.Command{
		Use:   "update_stats",
		Short: "Runs background process which updates creative stats",
		Run: func(c *cobra.Command, args []string) {
			cmd.Run()
		},
	}

	return &cmd
}

func init() {
	cmd := newStatsCmd()
	cobra.OnInitialize(cmd.Init)
	cmd.PersistentFlags().StringVar(&cmd.cfgFile, "config", "", "config file")
	rc.AddCommand(cmd.Command)
}
