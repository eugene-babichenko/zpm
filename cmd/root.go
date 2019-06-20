package cmd

import (
	"github.com/eugene-babichenko/zpm/config"
	"github.com/eugene-babichenko/zpm/meta"

	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gopkg.in/yaml.v2"
)

var appConfigFile string
var appConfig config.Config
var cachePath string
var metaFilePath string
var updateCheckPeriod time.Duration

var lastUpdate time.Time

var logger *zap.SugaredLogger

var RootCmd = &cobra.Command{
	Use:   "zpm [command]",
	Short: "A simple zsh plugin manager",
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println("failed to execute the command:", err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	RootCmd.PersistentFlags().StringVar(
		&appConfigFile,
		"config",
		"",
		fmt.Sprintf("Config file location (default: %s)", defaultRootPrompt),
	)
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("cannot access the home directory:", err)
		os.Exit(1)
	}

	if appConfigFile == "" {
		appConfigFile = filepath.Join(home, ".zpm.json")
	}

	configFile, err := ioutil.ReadFile(appConfigFile)
	if os.IsNotExist(err) {
		configData, err := json.MarshalIndent(config.DefaultConfig, "", "  ")
		if err != nil {
			fmt.Println("failed to write the default config:", err)
			os.Exit(1)
		}
		if err := ioutil.WriteFile(appConfigFile, configData, os.ModePerm); err != nil {
			fmt.Println("failed to write the default config:", err)
			os.Exit(1)
		}
	} else if err != nil {
		fmt.Println("failed to read the configuration file", err)
		os.Exit(1)
	} else {
		switch filepath.Ext(appConfigFile) {
		case "json":
			err = json.Unmarshal(configFile, &appConfig)
		case "yaml", "yml":
			err = yaml.Unmarshal(configFile, &appConfig)
		}
		if err != nil {
			fmt.Println("failed to parse the configuration file:", err)
			os.Exit(1)
		}
	}

	var logsPath string
	if len(appConfig.Root) == 0 {
		appConfig.Root = filepath.Join(home, defaultRoot)
		logsPath = filepath.Join(home, defaultLogs)
	} else {
		logsPath = filepath.Join(appConfig.Root, "Logs")
	}

	cachePath = filepath.Join(appConfig.Root, "cache.zsh")

	if appConfig.UpdateCheckPeriod != "" {
		updateCheckPeriodLocal, err := time.ParseDuration(appConfig.UpdateCheckPeriod)
		if err != nil {
			fmt.Println("failed to parse the update check period")
			os.Exit(1)
		}
		updateCheckPeriod = updateCheckPeriodLocal
	} else {
		updateCheckPeriodLocal, _ := time.ParseDuration("24h")
		updateCheckPeriod = updateCheckPeriodLocal
	}

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

	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = timeEncoder
	encoderConfig.EncodeLevel = levelEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(io.MultiWriter(os.Stdout, fileLogger)),
		level,
	)

	logger = zap.New(core).Sugar()

	metaFilePath = filepath.Join(appConfig.Root, "meta.json")
	metaFile, err := ioutil.ReadFile(metaFilePath)
	if err != nil {
		return
	}
	var metaData meta.Meta
	if err := json.Unmarshal(metaFile, &metaData); err != nil {
		return
	}
	lastUpdate, _ = time.Parse(meta.LastUpdateCheckLayout, metaData.LastUpdateCheck)
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("Mon Jan 2 15:04:05 2006"))
}

func levelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + l.CapitalString() + "]")
}
