package cmd

import (
	"github.com/eugene-babichenko/zpm/log"

	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install new plugins",
	Run: func(cmd *cobra.Command, args []string) {
		ps, err := makePluginStorage(appConfig.Root, appConfig.Plugins)
		if err != nil {
			log.Fatal("while reading plugin configurations: %s", err)
			return
		}

		ps.checkPluginInstalls()
		ps.installAll()
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
