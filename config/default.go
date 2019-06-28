package config

func DefaultConfig() Config {
	return Config{
		Plugins:           []string{},
		Root:              "",
		UpdateCheckPeriod: "24h",
		LogLevel:          "info",
	}
}
