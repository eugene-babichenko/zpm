package plugin

import (
	"io/ioutil"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Check the path creation
func TestGitRepoPath(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "")
	require.Empty(t, err, "cannot create temp dir")

	plugin := NewGit("github.com/username/repo", "master", tempDir)
	expectedPath := filepath.Join(tempDir, "github.com/username/repo")
	assert.Equal(t, expectedPath, plugin.Dir.Path, "wrong path")
}
