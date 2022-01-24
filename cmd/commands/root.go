package commands

import (
	"fmt"

	"github.com/pls87/creative-rotation/internal/config"
	"github.com/pls87/creative-rotation/internal/logger"
	"github.com/spf13/cobra"
)

var (
	cfg     config.Config
	cfgFile string
	logg    *logger.Logger

	rootCmd = &cobra.Command{
		Use:   "cr",
		Short: "",
		Long:  `<Some long desc here...>`,
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
