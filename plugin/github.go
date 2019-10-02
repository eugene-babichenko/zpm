// Note that this file uses `go-git`. This was considered more fast and reliable
// than using an external `git` binary.

package plugin

import (
	"github.com/pkg/errors"
	"path/filepath"
)

func MakeGitHub(root string, params map[string]string) (*Plugin, error) {
	username, usernamePrs := params["username"]
	if !usernamePrs {
		return nil, errors.New("missing username")
	}

	repo, repoPrs := params["repo"]
	if !repoPrs {
		return nil, errors.New("missing repo")
	}

	requiredRevision := params["version"]
	if requiredRevision == "" {
		requiredRevision = "master"
	}

	URL := filepath.Join("github.com", username, repo)
	git := NewGit(URL, requiredRevision, root)
	plugin := Plugin(&git)
	return &plugin, nil
}
