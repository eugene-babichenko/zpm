package plugin

import (
	"os"

	"github.com/pkg/errors"
)

type File struct {
	Path string
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

func (p File) GetPath() string {
	return p.Path
}
