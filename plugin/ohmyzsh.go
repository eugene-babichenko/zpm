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

// Other plugins can depend on Oh My Zsh so we need a single global instance of it.
var ohMyZshInstance *OhMyZsh

func MakeOhMyZsh(root string, params []string) (*Plugin, error) {
	if len(params) != 2 {
		return nil, errors.New("invalid number of parameters")
	}

	requiredRevision := "master"
	if params[1] != "" {
		requiredRevision = params[1]
	}

	if ohMyZshInstance != nil {
		plugin := Plugin(ohMyZshInstance)
		return &plugin, nil
	}

	URL := filepath.Join("github.com", "robbyrussell", "oh-my-zsh")
	git := NewGit(URL, requiredRevision, root)

	ohMyZshInstance = &OhMyZsh{git: git}
	plugin := Plugin(ohMyZshInstance)

	return &plugin, nil
}

func MakeOhMyZshPlugin(root string, params []string) (*Plugin, error) {
	if len(params) != 1 {
		return nil, errors.New("invalid number of parameters")
	}

	_, err := MakeOhMyZsh(root, []string{"", ""})
	if err != nil {
		return nil, errors.Wrap(err, "while instantiating Oh My Zsh")
	}

	ohMyZshPlugin := Plugin(ohMyZshInstance.LoadPlugin(params[0]))

	return &ohMyZshPlugin, nil
}

func MakeOhMyZshTheme(root string, params []string) (*Plugin, error) {
	if len(params) != 1 {
		return nil, errors.New("invalid number of parameters")
	}

	_, err := MakeOhMyZsh(root, []string{"", ""})
	if err != nil {
		return nil, errors.Wrap(err, "while instantiating Oh My Zsh")
	}

	ohMyZshTheme := Plugin(ohMyZshInstance.LoadTheme(params[0]))

	return &ohMyZshTheme, nil
}

func GetOhMyZsh() *Plugin {
	if ohMyZshInstance != nil {
		plugin := Plugin(ohMyZshInstance)
		return &plugin
	}
	return nil
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

func (p *OhMyZsh) CheckUpdate() (*string, error) {
	return p.git.CheckUpdate()
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

func (p *OhMyZsh) GetPath() *string {
	return p.git.GetPath()
}
