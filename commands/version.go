package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current zpm version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("zpm %s", Version)
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
