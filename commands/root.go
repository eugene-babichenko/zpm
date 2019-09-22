package commands

import (
	"github.com/eugene-babichenko/zpm/assets"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	updateAPILink = "https://api.github.com/repos/eugene-babichenko/zpm/releases/latest"
	updateLink    = "https://github.com/eugene-babichenko/zpm"
)

var (
	Version string

	appConfigFile     string
	rootDir           string
	pluginsSpecs      []string
	updateCheckPeriod time.Duration

	RootCmd = &cobra.Command{
		Use:   "zpm [command]",
		Short: "A simple zsh plugin manager",
	}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatalf("failed to execute the command: %s", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(
		&appConfigFile,
		"config",
		"",
		"Config file location (default: $HOME/.zpm.yaml)",
	)
}

// prefixedWriter allows to add "zsh: " between log lines
type prefixedWriter struct{}

func (prefixedWriter) Write(p []byte) (n int, err error) {
	// Writing logs to stderr is workaround. In `source <(zpm load)` the
	// `<(...)` consumes only what is written to stdout. Thus, writing logs to
	// stderr allows us to have nice logs while loading plugins.
	nPrefix, err := os.Stderr.Write([]byte("zpm: "))
	if err != nil {
		return nPrefix, err
	}
	np, err := os.Stderr.Write(p)
	return nPrefix + np, err
}

func initConfig() {
	formatter := &log.TextFormatter{}
	formatter.DisableLevelTruncation = true
	formatter.DisableTimestamp = true
	// this is required to have colored output with a custom writer
	formatter.ForceColors = true

	log.SetFormatter(formatter)
	log.SetOutput(prefixedWriter{})

	viper.SetConfigName(".zpm")
	viper.AddConfigPath("$HOME")

	viper.SetDefault("Plugins", []string{})
	viper.SetDefault("LoggingLevel", "info")
	viper.SetDefault("OnLoad.InstallMissingPlugins", true)
	viper.SetDefault("OnLoad.CheckForUpdates", true)
	viper.SetDefault("OnLoad.UpdateCheckPeriod", "24h")

	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get system root directory: %s", err)
	}
	rootDir = filepath.Join(home, ".zpm_plugins")

	if err := os.MkdirAll(rootDir, os.ModePerm); err != nil && !os.IsExist(err) {
		log.Fatalf("while creating the plugin storage directory: %s", err)
	}

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Fatalf("failed to read configuration: %s", err)
		}
		// write defaults
		configFilePath := filepath.Join(home, ".zpm.yaml")
		// this is used to preserve uppercase letters which viper does not care about
		yamlConfig, _ := assets.Asset("config/.zpm.yaml")
		if err := ioutil.WriteFile(configFilePath, yamlConfig, os.ModePerm); err != nil {
			log.Fatalf("failed to write the default config to the drive: %s", err)
		}
	}

	pluginsSpecs = viper.GetStringSlice("Plugins")

	level, err := log.ParseLevel(viper.GetString("LoggingLevel"))
	if err != nil {
		log.Errorf("failed to set the logging level: %s", err)
	}

	log.SetLevel(level)

	updateCheckPeriod, err = time.ParseDuration(viper.GetString("OnLoad.UpdateCheckPeriod"))
	if err != nil {
		log.Fatalf("failed to parse OnLoad.UpdateCheckPeriod")
	}
}
