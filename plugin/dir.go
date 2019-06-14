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

	entrypoint := filepath.Join(p.Path, filepath.Base(p.Path)+".plugin.zsh")
	stat, err = os.Stat(entrypoint)
	if err == nil {
		if stat.Mode()&os.ModeType == 0 {
			s = append(s, fmt.Sprintf("source %s", entrypoint))
		}
	}

	theme := filepath.Join(p.Path, filepath.Base(p.Path)+".zsh-theme")
	stat, err = os.Stat(theme)
	if err == nil {
		if stat.Mode()&os.ModeType == 0 {
			s = append(s, fmt.Sprintf("source %s", theme))
		}
	}

	return s, nil
}
