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

func (p Dir) Load() ([]string, error) {
	stat, err := os.Stat(p.Path)
	if err != nil {
		return nil, errors.Wrap(err, "while loading directory plugin")
	}
	if stat.Mode()&os.ModeType == 0 {
		return nil, errors.New("the provided path is not a directory: " + p.Path)
	}

	s := make([]string, 1)

	s[0] = fmt.Sprintf("fpath=($fpath %s)", p.Path)

	entrypoints, err := filepath.Glob(filepath.Join(p.Path, "*.plugin.zsh"))
	if err != nil {
		return nil, errors.Wrap(err, "while loading directory plugin")
	}

	themes, err := filepath.Glob(filepath.Join(p.Path, "*.zsh-theme"))
	if err != nil {
		return nil, errors.Wrap(err, "while loading directory plugin")
	}
	entrypoints = append(entrypoints, themes...)

	for _, entrypoint := range entrypoints {
		stat, err = os.Stat(entrypoint)
		if err == nil {
			if stat.Mode()&os.ModeType == 0 {
				s = append(s, fmt.Sprintf("source %s", entrypoint))
			}
		}
	}

	return s, nil
}

func (p Dir) CheckUpdate() (*string, error) {
	return nil, nil
}

func (p Dir) InstallUpdate() error {
	return nil
}
