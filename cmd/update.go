package cmd

import (
	"zpm/plugin"

	"os"
	"sync"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Install updates and download missing plugins",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("invalidating cache...")
		if err := os.RemoveAll(cachePath); err != nil {
			logger.Error("while invalidating cache: ", err.Error())
		}

		names, plugins, err := appConfig.GetPlugins()
		if err != nil {
			logger.Fatal("while reading plugin configurations: ", err.Error())
		}

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(plugins))

		for idx, pluginInstance := range plugins {
			go func(idx int, pluginInstance plugin.Plugin) {
				if update, err := pluginInstance.CheckUpdate(); plugin.IsNotInstalled(err) {
					logger.Info("installing: ", names[idx])
					if err := pluginInstance.InstallUpdate(); err != nil {
						logger.Errorf("installation error for %s: %s", names[idx], err.Error())
					}
					logger.Info("installed: ", names[idx])
				} else if err != nil {
					logger.Errorf("while checking for %s: %s", names[idx], err.Error())
				} else if update != nil {
					logger.Infof("updating %s: %s", names[idx], *update)
					if err := pluginInstance.InstallUpdate(); err != nil {
						logger.Errorf("while updating %s: %s", names[idx], err.Error())
					}
					logger.Info("updated: ", names[idx])
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
	RootCmd.AddCommand(updateCmd)
}
