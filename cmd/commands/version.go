package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Release   string
	BuildDate string
	GitHash   string
)

func init() {
	rc.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of Creative Rotation app",
	Run: func(cmd *cobra.Command, args []string) {
		if err := json.NewEncoder(os.Stdout).Encode(struct {
			Release   string
			BuildDate string
			GitHash   string
		}{
			Release:   Release,
			BuildDate: BuildDate,
			GitHash:   GitHash,
		}); err != nil {
			fmt.Printf("error while decode version info: %v\n", err)
		}
	},
}
