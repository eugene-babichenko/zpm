package config

func DefaultConfig() Config {
	return Config{
		Plugins:           []string{},
		Root:              "",
		UpdateCheckPeriod: "24h",
		LogsPath:          "",
		Logger: Logger{
			MaxSize:    500,
			MaxAge:     28,
			MaxBackups: 6,
			Level:      "info",
		},
	}
}
