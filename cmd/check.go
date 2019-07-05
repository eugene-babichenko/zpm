package cmd

import (
	"github.com/eugene-babichenko/zpm/log"
	"github.com/eugene-babichenko/zpm/plugin"

	"os"
	"sync"
	"sync/atomic"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for updates",
	Run: func(cmd *cobra.Command, args []string) {
		names, plugins, err := MakePluginsFromSpecs(appConfig.Root, appConfig.Plugins)
		if err != nil {
			log.Fatal("while reading plugin configurations: %s", err)
			os.Exit(1)
		}

		var updatesAvailable int32
		var installationsAvailable int32

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(plugins))

		for idx, pluginInstance := range plugins {
			go func(name string, pluginInstance plugin.Plugin) {
				if updateString, err := checkPluginUpdate(name, pluginInstance); updateString != nil {
					atomic.AddInt32(&updatesAvailable, 1)
				} else if plugin.IsNotInstalled(err) {
					atomic.AddInt32(&installationsAvailable, 1)
				}
				waitGroup.Done()
			}(names[idx], pluginInstance)
		}

		waitGroup.Wait()

		if updatesAvailable > 0 || installationsAvailable > 0 {
			log.Info(
				"%d updates available and %d plugins need to be installed",
				updatesAvailable,
				installationsAvailable,
			)
			log.Info("You can run the update using `zpm update`.")
		}

		updateLastUpdateCheckTime()
	},
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
