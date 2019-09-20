package cmd

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Install updates",
	Run: func(cmd *cobra.Command, args []string) {
		pluginToUpdate, _ := cmd.Flags().GetString("plugin")

		ps, err := makePluginStorage(rootDir, pluginsSpecs)
		if err != nil {
			log.Fatalf("while reading plugin configurations: %s", err)
		}

		if pluginToUpdate == "" {
			ps.checkPluginUpdates(false)
			ps.updateAll()
			if err := setLastUpdateTime(time.Now()); err != nil {
				log.Errorf("failed to write last update time: %s", err)
				log.Error("note that this will result in extra update checks on zsh load")
			}
			return
		}

		pse, ok := ps.plugins[pluginToUpdate]
		if !ok {
			log.Fatalf("plugin %s not listed in the configuration file", pluginToUpdate)
		}
		pse.checkPluginUpdate(false)
		pse.update()
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
