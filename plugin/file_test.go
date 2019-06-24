package plugin

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

// Feature: File plugin
//   Scenario: Successfully load a plugin
//     Given that a plugin file is in its place
//     When the `Load` function is called
//     Then the source line for this file is in the output
//     And `fpath` is empty
func TestFileLoadSuccess(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	pluginFileName := filepath.Join(tempDir, "test.plugin.zsh")
	_, err = os.Create(pluginFileName)
	require.Empty(t, err, "cannot create test file")

	pluginInstance, err := MakeFile(tempDir, map[string]string{"filename": "test.plugin.zsh"})
	require.Empty(t, err, "cannot create a plugin")

	fpath, exec, err := (*pluginInstance).Load()
	require.Empty(t, err, "cannot load an existing plugin")
	assert.Empty(t, fpath, "fpath must be empty for file plugins")
	assert.Equal(t, exec, []string{"source " + pluginFileName}, "incorrect exec line")
}

//   Scenario: Plugin not installed
//     Given that a plugin file does not exist
//     When the `Load` function is called
//     Then the error must be "plugin not installed"
func TestFileLoadNotExist(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	pluginInstance, err := MakeFile(tempDir, map[string]string{"filename": "test.plugin.zsh"})
	require.Empty(t, err, "cannot create a plugin")

	_, _, err = (*pluginInstance).Load()
	assert.Equal(t, NotInstalled, err, "the error must be 'plugin not installed'")
}
