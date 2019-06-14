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

func (c Config) GetPlugins() ([]plugin.Plugin, error) {
	plugins := make([]plugin.Plugin, 0, len(c.Plugins))

	for _, pluginSpec := range c.Plugins {
		submatch := filePluginRegex.FindStringSubmatch(pluginSpec)
		if len(submatch) > 0 {
			filename := submatch[1]
			plugins = append(plugins, plugin.File{Path: filepath.Join(c.Root, "plugins", filename)})
			continue
		}

		return nil, errors.New("unknown plugin format")
	}

	return plugins, nil
}
