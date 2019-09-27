package plugin

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"path/filepath"
	"testing"
)

// Feature: Oh My Zsh support
//   Scenario: Load plugin
func TestOhMyZshLoadPlugin(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	ohmyzsh, _ := MakeOhMyZsh(tempDir, map[string]string{})

	plugin := (*ohmyzsh).(*OhMyZsh).LoadPlugin("cargo")
	expectedPath := filepath.Join(tempDir, "github.com/robbyrussell/oh-my-zsh/plugins/cargo")
	assert.Equal(t, expectedPath, plugin.Path, "invalid plugin path")
}

//   Scenario: Load plugin (make function)
//     When the function is called with the required parameters
//     Then no error is returned
//     And a plugin instance is returned
func TestMakeOhMyZshPlugin(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	ohmyzsh, _ := MakeOhMyZsh(tempDir, map[string]string{})

	plugin, err := (*ohmyzsh).(*OhMyZsh).MakePlugin(tempDir, map[string]string{"name": "cargo"})
	assert.NotEmpty(t, plugin, "must return a plugin")
	assert.Empty(t, err, "must not return an error")
}

//   Scenario: Load plugin without name (make function)
//     When the function is called without the "name" argument
//     Then an error is returned
func TestMakeOhMyZshPluginNoName(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	ohmyzsh, _ := MakeOhMyZsh(tempDir, map[string]string{})

	_, err = (*ohmyzsh).(*OhMyZsh).MakePlugin(tempDir, map[string]string{})
	assert.NotEmpty(t, err, "must return an error")
}

//   Scenario: Load theme (make function)
//     When the function is called with the required parameters
//     Then no error is returned
//     And a plugin instance is returned
func TestMakeOhMyZshTheme(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	ohmyzsh, _ := MakeOhMyZsh(tempDir, map[string]string{})

	plugin, err := (*ohmyzsh).(*OhMyZsh).MakeTheme(tempDir, map[string]string{"name": "default"})
	assert.NotEmpty(t, plugin, "must return a plugin")
	assert.Empty(t, err, "must not return an error")
}

//   Scenario: Load theme without name (make function)
//     When the function is called without the "name" argument
//     Then an error is returned
func TestMakeOhMyZshThemeNoName(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	ohmyzsh, _ := MakeOhMyZsh(tempDir, map[string]string{})

	_, err = (*ohmyzsh).(*OhMyZsh).MakeTheme(tempDir, map[string]string{})
	assert.NotEmpty(t, err, "must return an error")
}
