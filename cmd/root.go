package cmd

import (
	"zpm/config"

	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var appConfigFile string
var appConfig config.Config
var cachePath string

var logger *zap.SugaredLogger

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
		fmt.Sprintf("config file (default: %s)", defaultRootPrompt),
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

	var logsPath string
	if len(appConfig.Root) == 0 {
		appConfig.Root = filepath.Join(home, defaultRoot)
		logsPath = filepath.Join(home, defaultLogs)
	} else {
		logsPath = filepath.Join(appConfig.Root, "Logs")
	}

	cachePath = filepath.Join(appConfig.Root, "cache.zsh")

	var level zapcore.Level
	switch appConfig.Logger.Level {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	fileLogger := &lumberjack.Logger{
		Filename:   filepath.Join(logsPath, "zpm.log"),
		MaxSize:    appConfig.Logger.MaxSize,
		MaxAge:     appConfig.Logger.MaxAge,
		MaxBackups: appConfig.Logger.MaxBackups,
		LocalTime:  true,
		Compress:   false,
	}
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig()),
		zapcore.AddSync(io.MultiWriter(os.Stdout, fileLogger)),
		level,
	)

	logger = zap.New(core).Sugar()
}
