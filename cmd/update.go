package cmd

import (
	"github.com/eugene-babichenko/zpm/log"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Install updates and download missing plugins",
	Run: func(cmd *cobra.Command, args []string) {
		onlyMissing, _ := cmd.Flags().GetBool("only-missing")
		pluginToCheck, _ := cmd.Flags().GetString("plugin")

		var pluginsList []string

		// Update a single plugin if required.
		if pluginToCheck != "" {
			pluginsList = []string{pluginToCheck}
		} else {
			pluginsList = appConfig.Plugins
		}

		names, plugins, err := MakePluginsFromSpecs(appConfig.Root, pluginsList)
		if err != nil {
			log.Fatal("while reading plugin configurations: %s", err)
		}

		checkAndInstallUpdates(names, plugins, !onlyMissing, true, false)
	},
}

func init() {
	updateCmd.Flags().String(
		"plugin",
		"",
		"Update only the specified plugin",
	)

	updateCmd.Flags().Bool(
		"only-missing",
		false,
		"Only install missing dependencies without updating the installed ones",
	)

	RootCmd.AddCommand(updateCmd)
}
