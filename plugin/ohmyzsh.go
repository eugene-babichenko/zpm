package plugin

import (
	"fmt"
	"path/filepath"

	"github.com/pkg/errors"
)

// The plugin type to deal with Oh My Zsh
type OhMyZsh struct {
	git Git
}

func MakeOhMyZsh(root string, params map[string]string) (*Plugin, error) {
	requiredRevision := params["version"]
	if requiredRevision == "" {
		requiredRevision = "master"
	}

	URL := filepath.Join("github.com", "robbyrussell", "oh-my-zsh")
	git := NewGit(URL, requiredRevision, root)

	omz := &OhMyZsh{git: git}
	plugin := Plugin(omz)

	return &plugin, nil
}

func (p OhMyZsh) MakePlugin(_ string, params map[string]string) (*Plugin, error) {
	plugin, pluginPrs := params["name"]
	if !pluginPrs {
		return nil, errors.New("missing plugin name")
	}

	ohMyZshPlugin := Plugin(p.LoadPlugin(plugin))

	return &ohMyZshPlugin, nil
}

func (p OhMyZsh) MakeTheme(_ string, params map[string]string) (*Plugin, error) {
	theme, themePrs := params["name"]
	if !themePrs {
		return nil, errors.New("missing theme name")
	}

	ohMyZshTheme := Plugin(p.LoadTheme(theme))

	return &ohMyZshTheme, nil
}

func (p *OhMyZsh) Load() (fpath []string, exec []string, err error) {
	fpath, exec, err = p.git.Load()
	if err != nil {
		return nil, nil, errors.Wrap(err, "ohmyzsh")
	}

	// load zsh library files
	libraries := fmt.Sprintf(
		"for config_file (%s/lib/*.zsh); do source $config_file; done",
		p.git.Dir.Path,
	)
	exec = append(exec, libraries)

	return fpath, exec, nil
}

func (p *OhMyZsh) CheckUpdate(offline bool) (*string, error) {
	return p.git.CheckUpdate(offline)
}

func (p *OhMyZsh) InstallUpdate() error {
	return p.git.InstallUpdate()
}

func (p *OhMyZsh) LoadPlugin(name string) Dir {
	path := filepath.Join(p.git.Dir.Path, "plugins", name)
	return Dir{Path: path}
}

func (p *OhMyZsh) LoadTheme(name string) Dir {
	path := filepath.Join(p.git.Dir.Path, "themes", name)
	return Dir{Path: path}
}

func (p *OhMyZsh) IsInstalled() (installed bool, err error) {
	return p.git.IsInstalled()
}
