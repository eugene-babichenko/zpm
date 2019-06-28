// The detailed documentation on those structs can be found in `README`.
// Explicit YAML annotations are required because go-yaml uses lowercase letters
// by default.
package config

import (
	"path/filepath"
	"strings"
)

type Config struct {
	Plugins           []string `yaml:"Plugins"`
	Root              string   `yaml:"Root"`
	UpdateCheckPeriod string   `yaml:"UpdateCheckPeriod"`
	LogLevel          string   `yaml:"LogLevel"`
}

// Validate assigns the default values to the config fields when applicable
func (c *Config) Validate(home string) {
	if c.Root == "" {
		c.Root = filepath.Join(home, DefaultRoot)
	}

	if c.UpdateCheckPeriod == "" {
		c.UpdateCheckPeriod = "24h"
	}

	if c.LogLevel == "" {
		c.LogLevel = "info"
	}

	c.LogLevel = strings.ToLower(c.LogLevel)
}
