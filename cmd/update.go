package cmd

import (
	"zpm/plugin"

	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Install updates and download missing plugins",
	Run: func(cmd *cobra.Command, args []string) {
		plugins, err := appConfig.GetPlugins()
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			os.Exit(1)
		}
		for idx, pluginInstance := range plugins {
			update, err := pluginInstance.CheckUpdate()
			if plugin.IsNotInstalled(err) {
				if err := pluginInstance.InstallUpdate(); err != nil {
					fmt.Printf("%s: installation error: %s\n", appConfig.Plugins[idx], err.Error())
				}
				continue
			} else if err != nil {
				fmt.Printf("%s: error: %s\n", appConfig.Plugins[idx], err.Error())
				continue
			}
			if update != nil {
				fmt.Printf("%s: updating: %s\n", appConfig.Plugins[idx], *update)
				if err := pluginInstance.InstallUpdate(); err != nil {
					fmt.Printf("%s: update error: %s\n", appConfig.Plugins[idx], err.Error())
				}
			} else {
				fmt.Printf("%s: up to date\n", appConfig.Plugins[idx])
			}
		}
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}