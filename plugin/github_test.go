package plugin

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test missing parameters
func TestMakeGitHubMissingUsername(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	pluginDir := filepath.Join(tempDir, "github.com/username/test")
	err = os.MkdirAll(pluginDir, os.ModePerm)
	require.Empty(t, err, "cannot create plugin dir")

	_, err = MakeGitHub(tempDir, map[string]string{"repo": "test"})
	assert.NotEmpty(t, err, "must return error")
}

func TestMakeGitHubMissingRepo(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	pluginDir := filepath.Join(tempDir, "github.com/username/test")
	err = os.MkdirAll(pluginDir, os.ModePerm)
	require.Empty(t, err, "cannot create plugin dir")

	_, err = MakeGitHub(tempDir, map[string]string{"username": "username"})
	assert.NotEmpty(t, err, "must return error")
}

func TestMakeGitHubSuccess(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	pluginDir := filepath.Join(tempDir, "github.com/username/test")
	err = os.MkdirAll(pluginDir, os.ModePerm)
	require.Empty(t, err, "cannot create plugin dir")

	_, err = MakeGitHub(tempDir, map[string]string{"username": "username", "repo": "test"})
	assert.Empty(t, err, "must not return error")
}
