package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Feature: Config default va;ue
//   Scenario: Logs path is empty
//     When the logs path is empty
//     And the validation is called
//     Then the logs path must be set to the default path
func TestConfigValidateLogsPath(t *testing.T) {
	config := DefaultConfig()
	config.Root = "/root"
	config.UpdateCheckPeriod = "30m"

	config.Validate("/home")
	assert.Equal(t, filepath.Join("/home", DefaultLogsPath), config.LogsPath, "unexpected logs path")
	assert.Equal(t, "/root", config.Root, "unexpected root path")
	assert.Equal(t, "30m", config.UpdateCheckPeriod, "unexpected update check period")
}

//   Scenario: Root is empty
//     When the root is empty
//     And the validation is called
//     Then the root must be set to the default path
func TestConfigValidateRoot(t *testing.T) {
	config := DefaultConfig()
	config.LogsPath = "/var/logs"
	config.UpdateCheckPeriod = "30m"

	config.Validate("/home")
	assert.Equal(t, "/var/logs", config.LogsPath, "unexpected logs path")
	assert.Equal(t, filepath.Join("/home", DefaultRoot), config.Root, "unexpected root path")
	assert.Equal(t, "30m", config.UpdateCheckPeriod, "unexpected update check period")
}

//   Scenario: Update period is empty
//     When the update period is empty
//     And the validation is called
//     Then the update period must be set to the default value
func TestConfigValidateUpdateCheckPeriod(t *testing.T) {
	config := DefaultConfig()
	config.LogsPath = "/var/logs"
	config.Root = "/home/.zpm"
	config.UpdateCheckPeriod = ""

	config.Validate("/home")
	assert.Equal(t, "/var/logs", config.LogsPath, "unexpected logs path")
	assert.Equal(t, "/home/.zpm", config.Root, "unexpected root path")
	assert.Equal(t, "24h", config.UpdateCheckPeriod, "unexpected update check period")
}
