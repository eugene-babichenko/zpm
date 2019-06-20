package cmd

import (
	"github.com/eugene-babichenko/zpm/plugin"

	"os"
	"sync"

	"github.com/spf13/cobra"
)

var onlyMissing bool

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Install updates and download missing plugins",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("invalidating cache...")
		if err := os.RemoveAll(cachePath); err != nil {
			logger.Error("while invalidating cache: ", err.Error())
		}

		names, plugins, err := MakePluginsFromSpecs(appConfig.Root, appConfig.Plugins)
		if err != nil {
			logger.Fatal("while reading plugin configurations: ", err.Error())
		}

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(plugins))

		for idx, pluginInstance := range plugins {
			go func(name string, pluginInstance plugin.Plugin) {
				defer waitGroup.Done()

				pluginPath := pluginInstance.GetPath()
				if pluginPath != nil {
					if stat, _ := os.Stat(*pluginPath); stat != nil && onlyMissing {
						return
					}
				}

				update, err := checkPluginUpdate(name, pluginInstance)

				if plugin.IsNotInstalled(err) {
					logger.Info("installing: ", name)
					if err := pluginInstance.InstallUpdate(); err != nil {
						logger.Errorf("installation error for %s: %s", name, err.Error())
					}
					logger.Info("installed: ", name)
				} else if err == nil && update != nil {
					logger.Infof("updating %s: %s", name, *update)
					if err := pluginInstance.InstallUpdate(); err != nil {
						logger.Errorf("while updating %s: %s", name, err.Error())
					}
					logger.Info("updated: ", names)
				}
			}(names[idx], pluginInstance)
		}

		waitGroup.Wait()
	},
}

func init() {
	updateCmd.Flags().BoolVar(
		&onlyMissing,
		"only-missing",
		false,
		"Only install missing dependencies without updating the installed ones",
	)

	RootCmd.AddCommand(updateCmd)
}
