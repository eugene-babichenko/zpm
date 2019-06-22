// Note that this file uses `go-git`. This was considered more fast and reliable
// than using an external `git` binary.

package plugin

import (
	"github.com/pkg/errors"
	"path/filepath"
)

func MakeGitHub(root string, params []string) (*Plugin, error) {
	if len(params) != 4 {
		return nil, errors.New("invalid number of parameters")
	}

	requiredRevision := "master"
	if params[3] != "" {
		requiredRevision = params[3]
	}

	URL := filepath.Join("github.com", params[0], params[1])
	git := NewGit(URL, requiredRevision, root)
	plugin := Plugin(&git)
	return &plugin, nil
}
