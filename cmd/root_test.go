package cmd

import (
	"github.com/eugene-babichenko/zpm/config"

	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

// Feature: Configuration
//   Scenario: The configuration file does not exist
//     Given that the config does not exist
//     When the config is requested
//     Then the config file is created
//     And filled with the default content
//     And the default content is returned
func TestDefaultConfigCreated(t *testing.T) {
	tests := []struct {
		extension string
		unmarshal func([]byte, interface{}) error
	}{
		{"yaml", yaml.Unmarshal},
		{"yml", yaml.Unmarshal},
		{"json", json.Unmarshal},
	}

	for _, tt := range tests {
		tempDir, err := ioutil.TempDir("", "")
		require.Empty(t, err, "cannot create temp dir")

		configFileName := filepath.Join(tempDir, ".zpm."+tt.extension)

		configVal, err := loadConfigOrCreateDefault(configFileName)
		require.Empty(t, err, "cannot load or create the config")
		assert.Equal(t, config.DefaultConfig(), *configVal, "not a default config returned")

		configFile, err := ioutil.ReadFile(configFileName)
		require.Empty(t, os.IsNotExist(err), "config not written")
		require.Empty(t, err, "cannot read a new config")

		var newConfigVal config.Config
		err = tt.unmarshal(configFile, &newConfigVal)
		require.Empty(t, err, "cannot parse a new config")
		assert.Equal(t, config.DefaultConfig(), newConfigVal, "not a default config in file")
	}
}

//   Scenario: The configuration file already exists
//     Given that the config file exist
//     When the config is requested
//     Then the correct content is returned
func TestConfigExists(t *testing.T) {
	tests := []struct {
		extension string
		marshal   func(interface{}) ([]byte, error)
	}{
		{"yaml", yaml.Marshal},
		{"yml", yaml.Marshal},
		{"json", json.Marshal},
	}

	for _, tt := range tests {
		tempDir, err := ioutil.TempDir("", "")
		require.Empty(t, err, "cannot create temp dir")

		configFileName := filepath.Join(tempDir, ".zpm."+tt.extension)

		expectedConfigVal := config.DefaultConfig()
		expectedConfigVal.Root = "/config"

		configContent, err := tt.marshal(expectedConfigVal)
		require.Empty(t, err, "failed to marshal the expected config")
		require.Empty(t, ioutil.WriteFile(configFileName, configContent, os.ModePerm), "cannot write the required config")

		configVal, err := loadConfigOrCreateDefault(configFileName)
		require.Empty(t, err, "cannot load or create the config")
		assert.Equal(t, expectedConfigVal, *configVal, "not an expected config returned")
	}
}

//   Scenario: Unsupported extension
//     Given that the config file path has an unsupported extension
//     And this file exists
//     When the config is requested
//     Then an error is returned
func TestConfigUnsupported(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")
	configFileName := filepath.Join(tempDir, ".zpm.jsn")
	_, err = os.Create(configFileName)
	require.Empty(t, err, "cannot create temp file")

	_, err = loadConfigOrCreateDefault(configFileName)
	require.NotEmpty(t, err, "expected error")
}
