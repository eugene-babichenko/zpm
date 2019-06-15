package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for updates",
	Run: func(cmd *cobra.Command, args []string) {
		plugins, err := appConfig.GetPlugins()
		if err != nil {
			fmt.Printf("%s", err.Error())
			os.Exit(1)
		}
		for idx, plugin := range plugins {
			update, err := plugin.CheckUpdate()
			if err != nil {
				fmt.Printf("%s: %s", appConfig.Plugins[idx], err.Error())
				continue
			}
			if update != nil {
				fmt.Printf("%s: %s", appConfig.Plugins[idx], update)
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
