package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install new plugins",
	Run: func(cmd *cobra.Command, args []string) {
		ps, err := makePluginStorage(rootDir, pluginsSpecs)
		if err != nil {
			log.Fatalf("while reading plugin configurations: %s", err)
		}

		ps.checkPluginInstalls()
		ps.installAll()
	},
}

func init() {
	RootCmd.AddCommand(installCmd)
}
