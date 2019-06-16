package cmd

import (
	"zpm/plugin"

	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Install updates and download missing plugins",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("invalidating cache...")
		if err := os.RemoveAll(cachePath); err != nil {
			fmt.Println("error invalidating cache:", err.Error())
		}

		names, plugins, err := appConfig.GetPlugins()
		if err != nil {
			fmt.Printf("error: %s\n", err.Error())
			os.Exit(1)
		}

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(plugins))

		for idx, pluginInstance := range plugins {
			go func(idx int, pluginInstance plugin.Plugin) {
				if update, err := pluginInstance.CheckUpdate(); plugin.IsNotInstalled(err) {
					fmt.Printf("%s: installing...\n", names[idx])
					if err := pluginInstance.InstallUpdate(); err != nil {
						fmt.Printf("%s: installation error: %s\n", names[idx], err.Error())
					}
					fmt.Printf("%s: installed\n", names[idx])
				} else if err != nil {
					fmt.Printf("%s: error: %s\n", names[idx], err.Error())
				} else if update != nil {
					fmt.Printf("%s: updating: %s\n", names[idx], *update)
					if err := pluginInstance.InstallUpdate(); err != nil {
						fmt.Printf("%s: update error: %s\n", names[idx], err.Error())
					}
					fmt.Printf("%s: updated\n", names[idx])
				} else {
					fmt.Printf("%s: up to date\n", names[idx])
				}
				waitGroup.Done()
			}(idx, pluginInstance)
		}

		waitGroup.Wait()
	},
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
