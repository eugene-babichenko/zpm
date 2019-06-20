package cmd

import (
	"github.com/eugene-babichenko/zpm/plugin"

	"os"
	"sync"

	"github.com/spf13/cobra"
)

var onlyMissing bool
var pluginToCheck string

func update(name string, pluginInstance plugin.Plugin) {
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
		logger.Info("updated: ", name)
	}
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Install updates and download missing plugins",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Debug("invalidating cache...")
		if err := os.RemoveAll(cachePath()); err != nil {
			logger.Error("while invalidating cache: ", err.Error())
		}

		if pluginToCheck != "" {
			var pluginFound bool
			for _, pluginSpec := range appConfig.Plugins {
				pluginFound = pluginFound || (pluginSpec == pluginToCheck)
			}

			if pluginFound {
				pluginInstance, err := plugin.MakePlugin(appConfig.Root, pluginToCheck)
				if err != nil {
					logger.Fatal("while reading plugin configuration: ", err.Error())
				}
				if pluginInstance == nil {
					logger.Fatal("cannot load plugin instance")
				}
				//nil not possible because the program will exit on `logger.Fatal`
				//noinspection GoNilness
				update(pluginToCheck, *pluginInstance)

				return
			}

			logger.Fatal("this plugin is not listed in the configuration")
		}

		names, plugins, err := MakePluginsFromSpecs(appConfig.Root, appConfig.Plugins)
		if err != nil {
			logger.Fatal("while reading plugin configurations: ", err.Error())
		}

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(plugins))

		for idx, pluginInstance := range plugins {
			go func(name string, pluginInstance plugin.Plugin) {
				update(name, pluginInstance)
				waitGroup.Done()
			}(names[idx], pluginInstance)
		}

		waitGroup.Wait()

		updateMeta()
	},
}

func init() {
	updateCmd.Flags().StringVar(
		&pluginToCheck,
		"plugin",
		"",
		"Update the specified plugin",
	)

	updateCmd.Flags().BoolVar(
		&onlyMissing,
		"only-missing",
		false,
		"Only install missing dependencies without updating the installed ones",
	)

	RootCmd.AddCommand(updateCmd)
}
