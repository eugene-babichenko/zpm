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
			go func(idx int, pluginInstance plugin.Plugin) {
				if update, err := pluginInstance.CheckUpdate(); plugin.IsNotInstalled(err) {
					logger.Info("not installed: ", names[idx])
				} else if err != nil {
					logger.Errorf("while checking for %s: %s", names[idx], err.Error())
				} else if update != nil {
					logger.Infof("update available for %s: %s", names[idx], *update)
				} else {
					logger.Info("up to date: ", names[idx])
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
