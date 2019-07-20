package cmd

import (
	"github.com/eugene-babichenko/zpm/log"

	"os"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for updates",
	Run: func(cmd *cobra.Command, args []string) {
		names, plugins, err := MakePluginsFromSpecs(appConfig.Root, appConfig.Plugins)
		if err != nil {
			log.Fatal("while reading plugin configurations: %s", err)
			os.Exit(1)
		}

		checkAndInstallUpdates(names, plugins, false, false)
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
