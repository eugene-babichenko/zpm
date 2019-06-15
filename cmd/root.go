package cmd

import (
	"zpm/config"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var appConfigFile string
var appConfig config.Config
var cachePath string

var RootCmd = &cobra.Command{
	Use:   "zpm [command]",
	Short: "A simple zsh plugin manager",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(
		&appConfigFile,
		"config",
		"",
		"config file (default: $HOME/.zpm.json)",
	)
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if appConfigFile == "" {
		appConfigFile = filepath.Join(home, ".zpm.json")
	}

	configFile, err := ioutil.ReadFile(appConfigFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := json.Unmarshal(configFile, &appConfig); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if len(appConfig.Root) == 0 {
		appConfig.Root = filepath.Join(home, ".zpm")
	}

	cachePath = filepath.Join(appConfig.Root, "cache.zsh")
}
