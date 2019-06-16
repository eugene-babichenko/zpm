package cmd

import (
	"zpm/plugin"

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
			logger.Fatal("while reading plugin configurations: ", err.Error())
			os.Exit(1)
		}

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(plugins))

		for idx, pluginInstance := range plugins {
			go func(name string, pluginInstance plugin.Plugin) {
				_, _ = checkPluginUpdate(name, pluginInstance)
				waitGroup.Done()
			}(names[idx], pluginInstance)
		}

		waitGroup.Wait()
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
