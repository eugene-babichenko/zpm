package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

//   Scenario: Root is empty
//     When the root is empty
//     And the validation is called
//     Then the root must be set to the default path
func TestConfigValidateRoot(t *testing.T) {
	config := DefaultConfig()
	config.UpdateCheckPeriod = "30m"
	config.LogLevel = "debug"

	config.Validate("/home")
	assert.Equal(t, filepath.Join("/home", DefaultRoot), config.Root, "unexpected root path")
	assert.Equal(t, "30m", config.UpdateCheckPeriod, "unexpected update check period")
	assert.Equal(t, "debug", config.LogLevel, "unexpected log level")
}

//   Scenario: Update period is empty
//     When the update period is empty
//     And the validation is called
//     Then the update period must be set to the default value
func TestConfigValidateUpdateCheckPeriod(t *testing.T) {
	config := DefaultConfig()
	config.Root = "/home/.zpm"
	config.UpdateCheckPeriod = ""
	config.LogLevel = "debug"

	config.Validate("/home")
	assert.Equal(t, "/home/.zpm", config.Root, "unexpected root path")
	assert.Equal(t, "24h", config.UpdateCheckPeriod, "unexpected update check period")
	assert.Equal(t, "debug", config.LogLevel, "unexpected log level")
}

//   Scenario: Log level is empty
//     When the log level is empty
//     And the validation is called
//     Then the log level must be set to INFO
func TestConfigValidateLogLevel(t *testing.T) {
	config := DefaultConfig()
	config.Root = "/home/.zpm"
	config.LogLevel = ""

	config.Validate("/home")
	assert.Equal(t, "/home/.zpm", config.Root, "unexpected root path")
	assert.Equal(t, "24h", config.UpdateCheckPeriod, "unexpected update check period")
	assert.Equal(t, "info", config.LogLevel, "unexpected log level")
}

//   Scenario: Log levels are case-insensitive
//     When the log level is empty
//     And the validation is called
//     Then the log level must be set to INFO
func TestConfigValidateLogLevelCase(t *testing.T) {
	config := DefaultConfig()
	config.Root = "/home/.zpm"
	config.LogLevel = "Debug"

	config.Validate("/home")
	assert.Equal(t, "/home/.zpm", config.Root, "unexpected root path")
	assert.Equal(t, "24h", config.UpdateCheckPeriod, "unexpected update check period")
	assert.Equal(t, "debug", config.LogLevel, "unexpected log level")
}
