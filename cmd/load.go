package cmd

import (
	"github.com/eugene-babichenko/zpm/log"
	"github.com/eugene-babichenko/zpm/plugin"

	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cobra"
)

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load configured plugins into the current shell",
	Run: func(cmd *cobra.Command, args []string) {
		names, plugins, err := MakePluginsFromSpecs(appConfig.Root, appConfig.Plugins)
		if err != nil {
			log.Error("cannot load plugins: %s", err)
			os.Exit(1)
		}

		updateCheck, _ := cmd.Flags().GetBool("update-check")
		installMissing, _ := cmd.Flags().GetBool("install-missing")

		updateCheckPeriod, err := time.ParseDuration(appConfig.UpdateCheckPeriod)
		if err != nil {
			log.Fatal("failed to parse the update check period")
		}

		shouldCheckUpdate := updateCheck && readLastUpdateCheckTime().Add(updateCheckPeriod).Before(time.Now())

		waitGroup := sync.WaitGroup{}
		waitGroup.Add(len(plugins))

		for idx, pluginInstance := range plugins {
			go func(name string, pluginInstance plugin.Plugin) {
				defer waitGroup.Done()

				_, _, err := pluginInstance.Load()
				shouldInstall := plugin.IsNotInstalled(err) && installMissing

				if !shouldInstall && !shouldCheckUpdate {
					return
				}

				update, err := pluginInstance.CheckUpdate()

				if plugin.IsNotInstalled(err) && installMissing {
					log.Info("installing: %s", name)
					if err := pluginInstance.InstallUpdate(); err != nil {
						log.Error("while installing %s: %s", name, err.Error())
						return
					}
					log.Info("installed: %s", name)
				} else if update != nil && shouldCheckUpdate {
					log.Info("update available for %s: %s", name, *update)
				} else if err != nil && err != plugin.NotUpgradable && err != plugin.UpToDate {
					log.Error("while checking for an update for %s: %s", name, err)
				}
			}(names[idx], pluginInstance)
		}

		if shouldCheckUpdate {
			updateLastUpdateCheckTime()
		}

		waitGroup.Wait()

		fpath := make([]string, 0)
		exec := make([]string, 0)

		lines := []string{
			// Use different compdump for different zsh versions (kindly
			// borrowed from Oh My Zsh).
			"ZSH_COMPDUMP=\"${ZDOTDIR:-${HOME}}/.zcompdump-${SHORT_HOST}-${ZSH_VERSION}\"",
			"autoload -U compaudit compinit",
		}

		for _, plugin := range plugins {
			fpathPlugin, execPlugin, err := plugin.Load()
			if err != nil {
				log.Error("while loading plugin: %s", err)
				continue
			}
			fpath = append(fpath, fpathPlugin...)
			exec = append(exec, execPlugin...)
		}

		var fpathBuilder strings.Builder
		_, _ = fmt.Fprint(&fpathBuilder, "fpath=(")
		for _, fpathEntry := range fpath {
			_, _ = fmt.Fprintf(&fpathBuilder, "%s ", fpathEntry)
		}
		_, _ = fmt.Fprint(&fpathBuilder, "$fpath)")

		// Load completions from `fpath`.
		lines = append(lines, fpathBuilder.String(), "compinit -u -C -d \"${ZSH_COMPDUMP}\"")
		lines = append(lines, exec...)

		for _, line := range lines {
			fmt.Println(line)
		}
	},
}

func init() {
	loadCmd.Flags().Bool(
		"update-check",
		false,
		"Check for updates once in a period defined in the settings (default: 24h).",
	)

	loadCmd.Flags().Bool(
		"install-missing",
		false,
		"Install plugins that are listed in the configuration but are not installed.",
	)

	RootCmd.AddCommand(loadCmd)
}
