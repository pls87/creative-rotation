package crcmd

import (
	"github.com/pls87/creative-rotation/cmd/cr/config"
	"github.com/spf13/cobra"
)

var (
	cfg     config.Config
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "cr",
		Short: "",
		Long:  `<Some long desc here...>`,
		Run: func(cmd *cobra.Command, args []string) {
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
}
