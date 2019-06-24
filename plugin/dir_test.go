package plugin

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Feature: Dir plugin
//   Scenario: Successfully load a plugin
//     Given that a directory is in its place
//     And a file with an extension *.plugin.zsh or *.zsh-theme is present
//     When the `Load` function is called
//     Then the a file (or files) is sourced
//     And the directory is added to fpath
func TestDirLoadSuccess(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	pluginDir := filepath.Join(tempDir, "plugin")
	err = os.MkdirAll(pluginDir, os.ModePerm)
	require.Empty(t, err, "cannot create plugin dir")

	files := []string{"hello.plugin.zsh", "world.plugin.zsh", "impretty.zsh-theme"}
	for _, filename := range files {
		_, err = os.Create(filepath.Join(pluginDir, filename))
		require.Empty(t, err, "cannot create plugin file")
	}

	plugin, err := MakeDir(tempDir, map[string]string{"directory": "plugin"})
	require.Empty(t, err, "cannot create a plugin object")

	fpath, exec, err := (*plugin).Load()
	require.Empty(t, err, "cannot load a valid plugin")
	assert.Equal(t, []string{pluginDir}, fpath, "invalid fpath")

	sourceLines := make([]string, 0)
	for _, filename := range files {
		sourceLines = append(sourceLines, "source "+filepath.Join(pluginDir, filename))
	}

	assert.Equal(t, sourceLines, exec, "invalid exec lines")
}

//   Scenario: Valid name but object is a directory
//     Given that a directory is in its place
//     And a directory with a name *.plugin.zsh or *.zsh-theme is present
//     When the `Load` function is called
//     Then exec lines are empty
//     And the directory is added to fpath
func TestDirLoadSuccessValidNameButNotFile(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	pluginDir := filepath.Join(tempDir, "plugin")
	err = os.MkdirAll(pluginDir, os.ModePerm)
	require.Empty(t, err, "cannot create plugin dir")

	files := []string{"hello.plugin.zsh", "world.plugin.zsh", "impretty.zsh-theme"}
	for _, filename := range files {
		err = os.MkdirAll(filepath.Join(pluginDir, filename), os.ModePerm)
		require.Empty(t, err, "cannot create plugin file")
	}

	plugin, err := MakeDir(tempDir, map[string]string{"directory": "plugin"})
	require.Empty(t, err, "cannot create a plugin object")

	fpath, exec, err := (*plugin).Load()
	require.Empty(t, err, "cannot load a valid plugin")
	assert.Equal(t, []string{pluginDir}, fpath, "invalid fpath")
	assert.Empty(t, exec, "invalid exec lines")
}

//   Scenario: Directory does not exist
//     Given that a directory does not exist
//     When the `Load` function is called
//     Then a "not installed" error is returned
func TestDirLoadNotExist(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	plugin, err := MakeDir(tempDir, map[string]string{"directory": "plugin"})
	require.Empty(t, err, "cannot create a plugin object")

	_, _, err = (*plugin).Load()
	assert.Equal(t, NotInstalled, err, "unexpected error")
}

//   Scenario: Object is not a directory
//     Given the provided path is not a directory
//     When the `Load` function is called
//     Then an error is returned
func TestDirLoadNotADirectory(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	_, err = os.Create(filepath.Join(tempDir, "plugin"))
	require.Empty(t, err, "cannot create a plugin file")

	plugin, err := MakeDir(tempDir, map[string]string{"directory": "plugin"})
	require.Empty(t, err, "cannot create a plugin object")

	_, _, err = (*plugin).Load()
	assert.NotEmpty(t, err, "expected error")
}
