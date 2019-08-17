package cmd

import (
	"github.com/eugene-babichenko/zpm/log"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install new plugins",
	Run: func(cmd *cobra.Command, args []string) {
		var pluginsList []string

		names, plugins, err := MakePluginsFromSpecs(appConfig.Root, pluginsList)
		if err != nil {
			log.Fatal("while reading plugin configurations: %s", err)
		}

		checkAndInstallUpdates(names, plugins, false, true, false)
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
