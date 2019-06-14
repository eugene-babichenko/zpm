package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "load",
	Short: "Load configured plugins into the current shell",
	Run: func(cmd *cobra.Command, args []string) {
		plugins, err := appConfig.GetPlugins()
		if err != nil {
			fmt.Printf("# %s", err.Error())
			os.Exit(1)
		}
		for _, plugin := range plugins {
			lines, err := plugin.Load()
			if err != nil {
				fmt.Printf("# %s", err.Error())
				continue
			}
			for _, line := range lines {
				fmt.Println(line)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(versionCmd)
}
