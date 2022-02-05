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
	retryGap = 2 * time.Second
)

var (
	cfg     config.Config
	cfgFile string
	logg    *logger.Logger

	rootCmd = &cobra.Command{
		Use:   "cr",
		Short: "Creative rotation app",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Noop. Exiting....")
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(beforeRun)
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file")
}

func beforeRun() {
	cfg = config.New(cfgFile)
	logg = logger.New(cfg.Log)
}
