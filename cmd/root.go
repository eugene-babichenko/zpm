package cmd

import (
	"github.com/eugene-babichenko/zpm/config"
	"github.com/pkg/errors"

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

var (
	appConfigFile string
	appConfig     config.Config

	logger *zap.SugaredLogger

	RootCmd = &cobra.Command{
		Use:   "zpm [command]",
		Short: "A simple zsh plugin manager",
	}
)

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
		fmt.Sprintf("Config file location (default: $HOME/.zpm.yaml)"),
	)
}

func initConfig() {
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println("cannot access the home directory:", err)
		os.Exit(1)
	}

	if appConfigFile == "" {
		appConfigFile = filepath.Join(home, ".zpm.yaml")
	}

	appConfigLocal, err := loadConfigOrCreateDefault(appConfigFile)
	if err != nil {
		fmt.Printf("failed to read the config: %s\n", err.Error())
		os.Exit(1)
	}
	appConfigLocal.Validate(home)
	appConfig = *appConfigLocal

	level, err := getLoggingLevel(appConfig.Logger.Level)
	if err != nil {
		fmt.Printf("failed to set the logging level: %s\n", err.Error())
	}

	fileLogger := &lumberjack.Logger{
		Filename:   filepath.Join(appConfig.LogsPath, "zpm.log"),
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
}

func getLoggingLevel(levelString string) (zapcore.Level, error) {
	var level zapcore.Level
	switch levelString {
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "error":
		level = zap.ErrorLevel
	case "fatal":
		level = zap.FatalLevel
	case "":
		level = zap.InfoLevel
	default:
		return level, errors.New("invalid logging level specification")
	}
	return level, nil
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("Mon Jan 2 15:04:05 2006"))
}

func levelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString("[" + l.CapitalString() + "]")
}

func loadConfigOrCreateDefault(path string) (*config.Config, error) {
	if configFile, err := ioutil.ReadFile(path); os.IsNotExist(err) {
		var configData []byte
		switch filepath.Ext(path) {
		case ".json":
			configData, err = json.Marshal(config.DefaultConfig)
		case ".yaml", ".yml":
			configData, err = yaml.Marshal(config.DefaultConfig)
		default:
			return nil, errors.New("unsupported extension")
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to build the config file contents")
		}
		if err := ioutil.WriteFile(path, configData, os.ModePerm); err != nil {
			return nil, errors.Wrap(err, "failed to write the config file")
		}
		return &config.DefaultConfig, nil
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
