package commands

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

const (
	updateLink = "https://github.com/eugene-babichenko/zpm"

	configKeyPlugins                     = "plugins"
	configKeyLoggingLevel                = "logging_level"
	configKeyOnLoadInstallMissingPlugins = "on_load.install_missing_plugins"
	configKeyOnLoadCheckForUpdates       = "on_load.check_for_updates"
	configKeyOnLoadUpdateCheckPeriod     = "on_load.update_check_period"
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

	viper.SetDefault(configKeyPlugins, []string{})
	viper.SetDefault(configKeyLoggingLevel, "info")
	viper.SetDefault(configKeyOnLoadInstallMissingPlugins, true)
	viper.SetDefault(configKeyOnLoadCheckForUpdates, true)
	viper.SetDefault(configKeyOnLoadUpdateCheckPeriod, "24h")

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
		allSettings := viper.AllSettings()
		allSettingsBytes, err := yaml.Marshal(allSettings)
		if err != nil {
			log.Fatalf("failed to serialize settings: %s", err)
		}
		if err := ioutil.WriteFile(configFilePath, allSettingsBytes, os.ModePerm); err != nil {
			log.Fatalf("failed to write the default config to the drive: %s", err)
		}
	}

	pluginsSpecs = viper.GetStringSlice(configKeyPlugins)

	level, err := log.ParseLevel(viper.GetString(configKeyLoggingLevel))
	if err != nil {
		log.Errorf("failed to set the logging level: %s", err)
	}

	log.SetLevel(level)

	updateCheckPeriod, err = time.ParseDuration(viper.GetString(configKeyOnLoadUpdateCheckPeriod))
	if err != nil {
		log.Fatalf("failed to parse OnLoad.UpdateCheckPeriod")
	}
}
