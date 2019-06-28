package cmd

import (
	"github.com/eugene-babichenko/zpm/config"
	"github.com/eugene-babichenko/zpm/log"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var (
	appConfigFile string
	appConfig     config.Config

	RootCmd = &cobra.Command{
		Use:   "zpm [command]",
		Short: "A simple zsh plugin manager",
	}
)

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		log.Fatal("failed to execute the command: %s", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(
		&appConfigFile,
		"config",
		"",
		fmt.Sprintf("Config file location (default: $HOME/.zpm.yaml)"),
	)
}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("cannot access the home directory: %s", err)
	}

	if appConfigFile == "" {
		appConfigFile = filepath.Join(home, ".zpm.yaml")
	}

	appConfigLocal, err := loadConfigOrCreateDefault(appConfigFile)
	if err != nil {
		log.Fatal("failed to read the config: %s", err)
	}

	//noinspection GoNilness
	appConfigLocal.Validate(home)
	//noinspection GoNilness
	appConfig = *appConfigLocal

	if err := os.MkdirAll(appConfig.Root, os.ModePerm); err != nil && !os.IsExist(err) {
		log.Fatal("while creating github plugin object: %s", err)
	}

	level, err := getLoggingLevel(appConfig.LogLevel)
	if err != nil {
		log.Error("failed to set the logging level: %s", err)
		return
	}

	log.SetLevel(level)
}

func getLoggingLevel(levelString string) (log.Level, error) {
	var level log.Level
	switch levelString {
	case "debug":
		level = log.DebugLevel
	case "info":
		level = log.InfoLevel
	case "error":
		level = log.ErrorLevel
	case "fatal":
		level = log.FatalLevel
	case "":
		level = log.InfoLevel
	default:
		return level, errors.New("invalid logging level specification")
	}
	return level, nil
}

func loadConfigOrCreateDefault(path string) (*config.Config, error) {
	if configFile, err := ioutil.ReadFile(path); os.IsNotExist(err) {
		var configData []byte
		switch filepath.Ext(path) {
		case ".json":
			configData, err = json.Marshal(config.DefaultConfig())
		case ".yaml", ".yml":
			configData, err = yaml.Marshal(config.DefaultConfig())
		default:
			return nil, errors.New("unsupported extension")
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to build the config file contents")
		}
		if err := ioutil.WriteFile(path, configData, os.ModePerm); err != nil {
			return nil, errors.Wrap(err, "failed to write the config file")
		}
		c := config.DefaultConfig()
		return &c, nil
	} else if err != nil {
		return nil, errors.Wrap(err, "failed to read the config file")
	} else {
		var configData config.Config
		switch filepath.Ext(path) {
		case ".json":
			err = json.Unmarshal(configFile, &configData)
		case ".yaml", ".yml":
			err = yaml.Unmarshal(configFile, &configData)
		default:
			return nil, errors.New("unsupported extension")
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse the config file")
		}
		return &configData, nil
	}
}
