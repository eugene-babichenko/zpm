package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func loadCache() bool {
	stat, err := os.Stat(cachePath)
	if err != nil {
		fmt.Println("# error reading cache", err.Error())
		return false
	}
	if stat.Mode()&os.ModeType != 0 {
		fmt.Println("# error reading cache: not a file")
		return false
	}
	fmt.Printf("source " + cachePath)
	return true
}

var versionCmd = &cobra.Command{
	Use:   "load",
	Short: "Load configured plugins into the current shell",
	Run: func(cmd *cobra.Command, args []string) {
		if loadCache() {
			return
		}

		plugins, err := appConfig.GetPlugins()
		if err != nil {
			fmt.Println("# cannot load plugins:", err.Error())
			os.Exit(1)
		}

		lines := make([]string, 0)

		for _, plugin := range plugins {
			linesPlugin, err := plugin.Load()
			if err != nil {
				fmt.Println("# error loading plugin:", err.Error())
				continue
			}
			lines = append(lines, linesPlugin...)
		}

		for _, line := range lines {
			fmt.Println(line)
		}

		cacheFile, err := os.Create(cachePath)
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
	RootCmd.AddCommand(versionCmd)
}
