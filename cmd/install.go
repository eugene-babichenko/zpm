package cmd

import (
	"github.com/eugene-babichenko/zpm/plugin"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install new plugins",
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("installing plugins...")
		log.Info("not updating plugins! Run `zpm update` to do it.")

		ps, err := plugin.MakePluginStorage(rootDir, pluginsSpecs)
		if err != nil {
			log.Fatalf("while reading plugin configurations: %s", err)
		}

		ps.CheckPluginInstalls()
		ps.InstallAll()

		log.Info("installation finished!")
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
