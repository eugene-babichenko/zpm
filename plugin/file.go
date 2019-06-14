package plugin

import (
	"os"

	"github.com/pkg/errors"
)

type File struct {
	Path string
}

func (p File) Load() ([]string, error) {
	stat, err := os.Stat(p.Path)
	if err != nil {
		return nil, errors.Wrap(err, "while loading file plugin")
	}

	if stat.Mode()&os.ModeType != 0 {
		return nil, errors.New("the provided path is not a file: " + p.Path)
	}

	s := make([]string, 1)
	s[0] = "source " + p.Path

	return s, nil
}
