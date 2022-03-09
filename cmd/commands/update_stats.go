package commands

import (
	"context"
	"os/signal"
	"sync"
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
}

func (sc *StatsCMD) consumeImpressions() (chan stats.Event, chan error, error) {
	impressions, errors, err := sc.stats.Consume("StatsUpdater", stats.ImpressionQueue)
	if err != nil {
		sc.logg.Errorf("couldn't consume impressions: %s", err)
		return nil, nil, err
	}

	return impressions, errors, nil
}

func (sc *StatsCMD) consumeConversions() (chan stats.Event, chan error, error) {
	conversions, errors, err := sc.stats.Consume("StatsUpdater", stats.ConversionQueue)
	if err != nil {
		sc.logg.Errorf("couldn't consume conversions: %s", err)
		return nil, nil, err
	}

	return conversions, errors, nil
}

func (sc *StatsCMD) waitForMessages(ctx context.Context, t string, events chan stats.Event, errors chan error,
	handler func(stats.Event) error) {
	var ev stats.Event
	var e error
	for ok := true; ok; {
		select {
		case e, ok = <-errors:
			if !ok {
				break
			}
			sc.logg.Errorf("error while consuming %s: %s", t, e)
		case ev, ok = <-events:
			if !ok {
				break
			}
			if e = handler(ev); e != nil {
				sc.logg.Errorf("couln't update %s stats by event %v: %s", t, ev, e)
			}
		case <-ctx.Done():
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
	defer sc.shutDown()

	sc.logg.Info("stats updater process starting...")
	defer func() {
		sc.logg.Info("stats updater process finished...")
	}()

	sc.storage = storage.New(sc.cfg.DB)
	sc.stats = stats.NewConsumer(sc.cfg.Queue)

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

	i, ie, e := sc.consumeImpressions()
	if e != nil {
		return
	}
	c, ce, e := sc.consumeConversions()
	if e != nil {
		return
	}

	group := sync.WaitGroup{}

	group.Add(1)
	go func() {
		defer group.Done()
		sc.waitForMessages(ctx, "impression", i, ie, func(e stats.Event) error {
			return sc.storage.Stats().TrackImpression(context.Background(), e.CreativeID, e.SlotID, e.SegmentID)
		})
	}()

	group.Add(1)
	go func() {
		defer group.Done()
		sc.waitForMessages(ctx, "conversion", c, ce, func(e stats.Event) error {
			return sc.storage.Stats().TrackConversion(context.Background(), e.CreativeID, e.SlotID, e.SegmentID)
		})
	}()

	group.Wait()
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
