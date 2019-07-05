package cmd

import (
	"github.com/eugene-babichenko/zpm/log"

	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func loadCache() bool {
	stat, err := os.Stat(cachePath())
	if err != nil {
		log.Debug("while reading cache %s", err)
		return false
	}
	if stat.Mode()&os.ModeType != 0 {
		log.Debug("while reading cache: not a file")
		return false
	}
	fmt.Printf("source " + cachePath())
	return true
}

var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load configured plugins into the current shell",
	Run: func(cmd *cobra.Command, args []string) {
		noCache, _ := cmd.Flags().GetBool("no-cache")

		if !noCache {
			if loadCache() {
				return
			}
		}

		_, plugins, err := MakePluginsFromSpecs(appConfig.Root, appConfig.Plugins)
		if err != nil {
			log.Error("cannot load plugins: %s", err)
			os.Exit(1)
		}

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

		if noCache {
			return
		}

		cacheFile, err := os.Create(cachePath())
		if err != nil {
			log.Fatal("cannot write cache: %s", err)
		}

		for _, line := range lines {
			if _, err := fmt.Fprintln(cacheFile, line); err != nil {
				log.Fatal("cannot write cache: %s", err)
			}
		}
	},
}

func init() {
	loadCmd.Flags().Bool(
		"no-cache",
		false,
		"Do not use and set cache when loading plugins",
	)

	RootCmd.AddCommand(loadCmd)
}
