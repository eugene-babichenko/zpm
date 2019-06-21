// The detailed documentation on those structs can be found in `README`.
// Explicit YAML annotations are required because go-yaml uses lowercase letters
// by default.
package config

type Logger struct {
	MaxSize    int    `yaml:"MaxSize"`
	MaxAge     int    `yaml:"MaxAge"`
	MaxBackups int    `yaml:"MaxBackups"`
	Level      string `yaml:"Level"`
}

type Config struct {
	Plugins           []string `yaml:"Plugins"`
	Root              string   `yaml:"Root"`
	Logger            Logger   `yaml:"Logger"`
	UpdateCheckPeriod string   `yaml:"UpdateCheckPeriod"`
	LogsPath          string   `yaml:"LogsPath"`
}

// Validate assigns the default values to the config fields when applicable
func (c *Config) Validate() {
	if c.Root == "" {
		c.Root = DefaultRoot
	}
	if c.LogsPath == "" {
		c.LogsPath = DefaultLogsPath
	}
	if c.UpdateCheckPeriod == "" {
		c.UpdateCheckPeriod = "24h"
	}
}
