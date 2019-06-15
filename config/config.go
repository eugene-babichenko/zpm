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

func (c Config) GetPlugins() ([]plugin.Plugin, error) {
	plugins := make([]plugin.Plugin, 0, len(c.Plugins))

	for _, pluginSpec := range c.Plugins {
		if submatch := filePluginRegex.FindStringSubmatch(pluginSpec); len(submatch) > 0 {
			filename := submatch[1]
			plugins = append(plugins, plugin.File{Path: filepath.Join(c.Root, "plugins", filename)})
			continue
		}

		if submatch := dirPluginRegex.FindStringSubmatch(pluginSpec); len(submatch) > 0 {
			filename := submatch[1]
			plugins = append(plugins, plugin.Dir{Path: filepath.Join(c.Root, "plugins", filename)})
			continue
		}

		if submatch := githubPluginRegex.FindStringSubmatch(pluginSpec); len(submatch) > 0 {
			username := submatch[1]
			repositoryName := submatch[2]

			githubPlugin, err := plugin.NewGitHub(username, repositoryName, "branch", "master", c.Root)
			if err != nil {
				return nil, errors.Wrap(err, "failed to open a repository")
			}
			plugins = append(plugins, githubPlugin)
			continue
		}

		return nil, errors.New("unknown plugin format")
	}

	return plugins, nil
}
