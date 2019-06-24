package plugin

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// The plugin type loaded from a file. Not that id does not support `fpath`.
type File struct {
	Path string
}

func MakeFile(root string, params map[string]string) (*Plugin, error) {
	path, pathPrs := params["directory"]
	if !pathPrs {
		return nil, errors.New("missing path")
	}

	plugin := Plugin(File{Path: filepath.Join(root, path)})

	return &plugin, nil
}

func (p File) Load() (fpath []string, exec []string, err error) {
	stat, err := os.Stat(p.Path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "while loading file plugin")
	}

	if stat.Mode()&os.ModeType != 0 {
		return nil, nil, errors.New("the provided path is not a file: " + p.Path)
	}

	exec = []string{"source " + p.Path}

	return fpath, exec, err
}

func (p File) CheckUpdate() (*string, error) {
	return nil, nil
}

func (p File) InstallUpdate() error {
	return nil
}

func (p File) GetPath() *string {
	return &p.Path
}
