package plugin

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
)

type OhMyZsh struct {
	root   string
	github *GitHub
}

func NewOhMyZsh(root string) (*OhMyZsh, error) {
	github, err := NewGitHub("robbyrussell", "oh-my-zsh", "branch", "master", root)
	if err != nil {
		return nil, errors.Wrap(err, "ohmyzsh")
	}

	ohMyZsh := &OhMyZsh{
		root:   filepath.Join(root, "plugins", "github.com", "robbyrussell", "oh-my-zsh"),
		github: github,
	}

	return ohMyZsh, nil
}

func (p *OhMyZsh) Load() ([]string, error) {
	lines, err := p.github.Load()
	if err != nil {
		return nil, errors.Wrap(err, "ohmyzsh")
	}

	libraries := fmt.Sprintf(
		"for config_file (%s/lib/*.zsh); do source $config_file; done",
		p.github.Dir.Path,
	)

	lines = append(lines, "compinit -u -C", libraries)

	return lines, nil
}

func (p *OhMyZsh) CheckUpdate() (*string, error) {
	return p.github.CheckUpdate()
}

func (p *OhMyZsh) InstallUpdate() error {
	return p.github.InstallUpdate()
}

func (p *OhMyZsh) LoadPlugin(name string) Dir {
	path := filepath.Join(p.root, "plugins", name)
	return Dir{Path: path}
}

func (p *OhMyZsh) LoadTheme(name string) Dir {
	path := filepath.Join(p.root, "themes", name)
	return Dir{Path: path}
}
