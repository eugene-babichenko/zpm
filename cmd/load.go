package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var noCache bool

func loadCache() bool {
	if noCache {
		return false
	}

	stat, err := os.Stat(cachePath())
	if err != nil {
		fmt.Println("# error reading cache", err.Error())
		return false
	}
	if stat.Mode()&os.ModeType != 0 {
		fmt.Println("# error reading cache: not a file")
		return false
	}
	fmt.Printf("source " + cachePath())
	return true
}

var versionCmd = &cobra.Command{
	Use:   "load",
	Short: "Load configured plugins into the current shell",
	Run: func(cmd *cobra.Command, args []string) {
		if loadCache() {
			return
		}

		_, plugins, err := MakePluginsFromSpecs(appConfig.Root, appConfig.Plugins)
		if err != nil {
			fmt.Println("# cannot load plugins:", err.Error())
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
				fmt.Println("# error loading plugin:", err.Error())
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
			fmt.Println("# cannot write cache:", err.Error())
			os.Exit(1)
		}

		for _, line := range lines {
			if _, err := fmt.Fprintln(cacheFile, line); err != nil {
				fmt.Println("# cannot write cache:", err.Error())
				os.Exit(1)
			}
		}
	},
}

func init() {
	versionCmd.Flags().BoolVar(
		&noCache,
		"no-cache",
		false,
		"Do not use and set cache when loading plugins",
	)

	RootCmd.AddCommand(versionCmd)
}
