package cmd

import (
	"zpm/plugin"

	"fmt"
	"os"
	"sync"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for updates",
	Run: func(cmd *cobra.Command, args []string) {
		names, plugins, err := appConfig.GetPlugins()
		if err != nil {
			fmt.Printf("%s", err.Error())
			os.Exit(1)
		}

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(plugins))

		for idx, pluginInstance := range plugins {
			go func(idx int, pluginInstance plugin.Plugin) {
				if update, err := pluginInstance.CheckUpdate(); err != nil {
					fmt.Printf("%s: error: %s\n", names[idx], err.Error())
				} else if update != nil {
					fmt.Printf("%s: %s\n", names[idx], *update)
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
	RootCmd.AddCommand(checkCmd)
}
