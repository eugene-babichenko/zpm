package plugin

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// Dir is the plugin type loaded from a source directory.
type Dir struct {
	Path string
}

func MakeDir(root string, params map[string]string) (*Plugin, error) {
	path, pathPrs := params["directory"]
	if !pathPrs {
		return nil, errors.New("missing path")
	}

	plugin := Plugin(Dir{Path: filepath.Join(root, path)})

	return &plugin, nil
}

func (p Dir) Load() (fpath []string, exec []string, err error) {
	stat, err := os.Stat(p.Path)
	if err != nil {
		return nil, nil, NotInstalled
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
		stat, err = os.Stat(entrypoint)
		if err != nil {
			continue
		}
		if stat.Mode()&os.ModeType == 0 {
			exec = append(exec, fmt.Sprintf("source %s", entrypoint))
		}
	}

	return fpath, exec, nil
}

func (p Dir) CheckUpdate(bool) (*string, error) {
	return nil, ErrNotUpgradable
}

func (p Dir) InstallUpdate() error {
	return ErrNotUpgradable
}

func (p Dir) IsInstalled() (installed bool, err error) {
	return false, NotInstallable
}
