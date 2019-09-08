package cmd

import (
	"github.com/eugene-babichenko/zpm/log"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for updates",
	Run: func(cmd *cobra.Command, args []string) {
		ps, err := makePluginStorage(rootDir, pluginsSpecs)
		if err != nil {
			log.Fatal("while reading plugin configurations: %s", err)
			return
		}

		ps.checkPluginUpdates()
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
