package plugin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

type Dir struct {
	Path string
}

func MakeDir(root string, params []string) (*Plugin, error) {
	if len(params) != 1 {
		return nil, errors.New("invalid number of parameters")
	}

	plugin := Plugin(Dir{Path: filepath.Join(root, params[0])})

	return &plugin, nil
}

func (p Dir) Load() (fpath []string, exec []string, err error) {
	stat, err := os.Stat(p.Path)
	if err != nil {
		return nil, nil, errors.Wrap(err, "while loading directory plugin")
	}
	if stat.Mode()&os.ModeType == 0 {
		return nil, nil, errors.New("the provided path is not a directory: " + p.Path)
	}

	fpath = []string{fmt.Sprintf(p.Path)}

	entrypoints, err := filepath.Glob(filepath.Join(p.Path, "*.plugin.zsh"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "while loading directory plugin")
	}

	themes, err := filepath.Glob(filepath.Join(p.Path, "*.zsh-theme"))
	if err != nil {
		return nil, nil, errors.Wrap(err, "while loading directory plugin")
	}
	entrypoints = append(entrypoints, themes...)

	for _, entrypoint := range entrypoints {
		if stat, err = os.Stat(entrypoint); err == nil {
			if stat.Mode()&os.ModeType == 0 {
				exec = append(exec, fmt.Sprintf("source %s", entrypoint))
			}
		}
	}

	return fpath, exec, nil
}

func (p Dir) CheckUpdate() (*string, error) {
	return nil, nil
}

func (p Dir) InstallUpdate() error {
	return nil
}

func (p Dir) GetPath() *string {
	return &p.Path
}
