package cmd

import (
	"github.com/eugene-babichenko/zpm/plugin"

	"os"
	"sync"
	"sync/atomic"
	"time"

	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check for updates",
	Long: `
Check if there are any updates available.

Adding --periodic will cause the command to check for updates only after a
period defined in the settings have passed since the last check.
`,
	Run: func(cmd *cobra.Command, args []string) {
		periodicCheck, _ := cmd.Flags().GetBool("periodic")

		if periodicCheck {
			updateCheckPeriod, err := time.ParseDuration(appConfig.UpdateCheckPeriod)
			if err != nil {
				logger.Fatal("failed to parse the update check period")
			}

			if readLastUpdateCheckTime().Add(updateCheckPeriod).After(time.Now()) {
				return
			}
		}

		names, plugins, err := MakePluginsFromSpecs(appConfig.Root, appConfig.Plugins)
		if err != nil {
			logger.Fatal("while reading plugin configurations: ", err.Error())
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
			logger.Infof(
				"%d updates available and %d plugins need to be installed",
				updatesAvailable,
				installationsAvailable,
			)
			logger.Info("You can run the update using `zpm update`.")
		}

		updateLastUpdateCheckTime()
	},
}

func init() {
	checkCmd.Flags().Bool(
		"periodic",
		false,
		"Check only once in a period defined in the settings (default: 24h)",
	)

	RootCmd.AddCommand(checkCmd)
}
