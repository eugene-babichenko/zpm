package config

import (
	"github.com/pkg/errors"
	"zpm/plugin"
)

type Logger struct {
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Level      string
}

type Config struct {
	Plugins []string
	Root    string
	Logger  Logger
}

func (c Config) GetPlugins() (names []string, plugins []plugin.Plugin, err error) {
	for _, pluginSpec := range c.Plugins {
		p, err := plugin.MakePlugin(c.Root, pluginSpec)
		if err != nil {
			return nil, nil, errors.Wrap(err, "while loading plugins")
		}
		if pluginSpec != "oh-my-zsh" {
			plugins = append(plugins, *p)
			names = append(names, pluginSpec)
		}
	}

	if ohMyZsh := plugin.GetOhMyZsh(); ohMyZsh != nil {
		plugins = append([]plugin.Plugin{*ohMyZsh}, plugins...)
		names = append([]string{"oh-my-zsh"}, names...)
	}

	return names, plugins, nil
}
