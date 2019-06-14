package plugin

import (
	"os"

	"github.com/pkg/errors"
)

type File struct {
	Path string
}

func (p File) Load() ([]string, error) {
	if _, err := os.Stat(p.Path); err != nil {
		return nil, errors.Wrap(err, "while loading file plugin")
	}

	s := make([]string, 1)
	s[0] = "source " + p.Path

	return s, nil
}
