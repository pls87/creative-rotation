package commands

import (
	"fmt"
	"time"

	"github.com/pls87/creative-rotation/internal/config"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/spf13/cobra"
)

const (
	retries  = 5
	retryGap = 5 * time.Second
)

var rc *RootCMD

type RootCMD struct {
	*cobra.Command
	cfgFile string        // nolint: structcheck
	cfg     config.Config // nolint: structcheck
	logg    *logger.Logger
}

func (rc *RootCMD) Run() {
	fmt.Println("Noop. Exiting...")
}

func (rc *RootCMD) Retry(toRetry func() error, onError func()) {
	var err error
	for r := retries; r > 0; r-- {
		if err = toRetry(); err == nil {
			break
		}
		rc.logg.Errorf("failed to connect: %s", err)
		rc.logg.Info("retrying...")
		time.Sleep(retryGap)
	}

	if err != nil {
		rc.logg.Errorf("number of retries exceeded: %s", err)
		onError()
	}
}

func NewRootCommand() *RootCMD {
	cmd := new(RootCMD)
	cmd.Command = &cobra.Command{
		Use:   "cr",
		Short: "Creative rotation app",
		Run: func(c *cobra.Command, args []string) {
			cmd.Run()
		},
	}
	return cmd
}

func Execute() error {
	return rc.Execute()
}

func init() {
	rc = NewRootCommand()
}
