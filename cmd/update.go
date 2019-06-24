package cmd

import (
	"github.com/eugene-babichenko/zpm/plugin"

	"os"
	"sync"

	"github.com/spf13/cobra"
)

var (
	onlyMissing   bool
	pluginToCheck string
)

func update(name string, pluginInstance plugin.Plugin, onlyMissing bool) {
	update, err := checkPluginUpdate(name, pluginInstance)

	if plugin.IsNotInstalled(err) {
		logger.Info("installing: ", name)
		if err := pluginInstance.InstallUpdate(); err != nil {
			logger.Errorf("installation error for %s: %s", name, err.Error())
			return
		}
		logger.Info("installed: ", name)
	} else if err == nil && update != nil && !onlyMissing {
		logger.Infof("updating %s: %s", name, *update)
		if err := pluginInstance.InstallUpdate(); err != nil {
			logger.Errorf("while updating %s: %s", name, err.Error())
			return
		}
		logger.Info("updated: ", name)
	} else if err != nil {
		logger.Errorf("error while checking for an update: %s", err)
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

		var pluginsList []string

		// Update a single plugin if required.
		if pluginToCheck != "" {
			pluginsList = []string{pluginToCheck}
		} else {
			pluginsList = appConfig.Plugins
		}

		names, plugins, err := MakePluginsFromSpecs(appConfig.Root, pluginsList)
		if err != nil {
			logger.Fatal("while reading plugin configurations: ", err.Error())
		}

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(plugins))

		for idx, pluginInstance := range plugins {
			go func(name string, pluginInstance plugin.Plugin) {
				update(name, pluginInstance, onlyMissing)
				waitGroup.Done()
			}(names[idx], pluginInstance)
		}

		waitGroup.Wait()

		updateLastUpdateCheckTime()
	},
}

func init() {
	updateCmd.Flags().StringVar(
		&pluginToCheck,
		"plugin",
		"",
		"Update only the specified plugin",
	)

	updateCmd.Flags().BoolVar(
		&onlyMissing,
		"only-missing",
		false,
		"Only install missing dependencies without updating the installed ones",
	)

	RootCmd.AddCommand(updateCmd)
}
