package cmd

import (
	"github.com/eugene-babichenko/zpm/plugin"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for updates",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("checking for updates...")

		ps, err := plugin.MakePluginStorage(rootDir, pluginsSpecs)
		if err != nil {
			log.Fatalf("while reading plugin configurations: %s", err)
		}

		ps.CheckPluginUpdates(false)

		log.Info("update check finished")
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
