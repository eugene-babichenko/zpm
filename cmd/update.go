package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Install updates",
	Run: func(cmd *cobra.Command, args []string) {
		pluginToCheck, _ := cmd.Flags().GetString("plugin")

		var pluginsList []string
		// Update a single plugin if required.
		if pluginToCheck != "" {
			pluginsList = []string{pluginToCheck}
		} else {
			pluginsList = pluginsSpecs
		}

		ps, err := makePluginStorage(rootDir, pluginsList)
		if err != nil {
			log.Fatalf("while reading plugin configurations: %s", err)
		}

		ps.checkPluginUpdates()
		ps.updateAll()
	},
}

func init() {
	updateCmd.Flags().String(
		"plugin",
		"",
		"Update only the specified plugin",
	)

	RootCmd.AddCommand(updateCmd)
}
