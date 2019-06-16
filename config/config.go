package config

import (
	"zpm/plugin"

	"path/filepath"
	"regexp"

	"github.com/pkg/errors"
)

type Config struct {
	Plugins []string
	Root    string
}

var filePluginRegex = regexp.MustCompile(`file:(.*)`)
var dirPluginRegex = regexp.MustCompile(`dir:(.*)`)
var githubPluginRegex = regexp.MustCompile(`github:([a-z0-9\-]+)/([a-z0-9\-]+)`)
var ohMyZshPluginRegex = regexp.MustCompile(`oh-my-zsh:plugin:([a-z0-9\-]+)`)
var ohMyZshThemeRegex = regexp.MustCompile(`oh-my-zsh:theme:([a-z0-9\-]+)`)

func (c Config) GetPlugins() (names []string, plugins []plugin.Plugin, err error) {
	var ohMyZsh *plugin.OhMyZsh
	ohMyZshPlugins := make([]plugin.Plugin, 0)
	ohMyZshNames := make([]string, 0)

	loadOhMyZsh := func() error {
		if ohMyZsh == nil {
			ohMyZshLocal, err := plugin.NewOhMyZsh(c.Root)
			if err != nil {
				return errors.Wrap(err, "failed to open a repository")
			}

			ohMyZsh = ohMyZshLocal
		}

		return nil
	}

	for _, pluginSpec := range c.Plugins {
		if submatch := filePluginRegex.FindStringSubmatch(pluginSpec); len(submatch) > 0 {
			filename := submatch[1]
			plugins = append(plugins, plugin.File{Path: filepath.Join(c.Root, "plugins", filename)})
			names = append(names, pluginSpec)
			continue
		}

		if submatch := dirPluginRegex.FindStringSubmatch(pluginSpec); len(submatch) > 0 {
			filename := submatch[1]
			plugins = append(plugins, plugin.Dir{Path: filepath.Join(c.Root, "plugins", filename)})
			names = append(names, pluginSpec)
			continue
		}

		if submatch := githubPluginRegex.FindStringSubmatch(pluginSpec); len(submatch) > 0 {
			username := submatch[1]
			repositoryName := submatch[2]

			githubPlugin, err := plugin.NewGitHub(username, repositoryName, "branch", "master", c.Root)
			if err != nil {
				return nil, nil, errors.Wrap(err, "failed to open a repository")
			}
			plugins = append(plugins, githubPlugin)
			names = append(names, pluginSpec)
			continue
		}

		if pluginSpec == "ohmyzsh" {
			if err := loadOhMyZsh(); err != nil {
				return nil, nil, err
			}
			continue
		}

		if submatch := ohMyZshPluginRegex.FindStringSubmatch(pluginSpec); len(submatch) > 0 {
			if err := loadOhMyZsh(); err != nil {
				return nil, nil, err
			}
			ohMyZshPlugins = append(ohMyZshPlugins, ohMyZsh.LoadPlugin(submatch[1]))
			ohMyZshNames = append(ohMyZshNames, pluginSpec)
			continue
		}

		if submatch := ohMyZshThemeRegex.FindStringSubmatch(pluginSpec); len(submatch) > 0 {
			if err := loadOhMyZsh(); err != nil {
				return nil, nil, err
			}
			ohMyZshPlugins = append(ohMyZshPlugins, ohMyZsh.LoadTheme(submatch[1]))
			ohMyZshNames = append(ohMyZshNames, pluginSpec)
			continue
		}

		return nil, nil, errors.New("unknown plugin format")
	}

	if ohMyZsh != nil {
		ohMyZshPlugins = append(ohMyZshPlugins, ohMyZsh)
		plugins = append(ohMyZshPlugins, plugins...)

		ohMyZshNames = append(ohMyZshNames, "oh-my-zsh")
		names = append(ohMyZshNames, names...)
	}

	return names, plugins, nil
}
