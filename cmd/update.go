package cmd

import (
	"github.com/eugene-babichenko/zpm/log"
	"github.com/eugene-babichenko/zpm/plugin"

	"os"
	"sync"

	"github.com/spf13/cobra"
)

func update(name string, pluginInstance plugin.Plugin, onlyMissing bool) {
	if _, _, err := pluginInstance.Load(); !(plugin.IsNotInstalled(err) && onlyMissing) {
		return
	}

	update, err := checkPluginUpdate(name, pluginInstance)

	if plugin.IsNotInstalled(err) {
		log.Info("installing: %s", name)
		if err := pluginInstance.InstallUpdate(); err != nil {
			log.Error("installation error for %s: %s", name, err.Error())
			return
		}
		log.Info("installed: %s", name)
	} else if err == nil && update != nil && !onlyMissing {
		log.Info("updating %s: %s", name, *update)
		if err := pluginInstance.InstallUpdate(); err != nil {
			log.Error("while updating %s: %s", name, err)
			return
		}
		log.Info("updated: %s", name)
	} else if err != nil && err != plugin.NotUpgradable && err != plugin.UpToDate {
		log.Error("error while checking for an update: %s", err)
	}
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Install updates and download missing plugins",
	Run: func(cmd *cobra.Command, args []string) {
		onlyMissing, _ := cmd.Flags().GetBool("only-missing")
		pluginToCheck, _ := cmd.Flags().GetString("plugin")

		log.Debug("invalidating cache...")
		if err := os.RemoveAll(cachePath()); err != nil {
			log.Error("while invalidating cache: %s", err)
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
			log.Fatal("while reading plugin configurations: %s", err)
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
	updateCmd.Flags().String(
		"plugin",
		"",
		"Update only the specified plugin",
	)

	updateCmd.Flags().Bool(
		"only-missing",
		false,
		"Only install missing dependencies without updating the installed ones",
	)

	RootCmd.AddCommand(updateCmd)
}
