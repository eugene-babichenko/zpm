package plugin

import (
	"regexp"

	"github.com/pkg/errors"
)

var ErrUnknownPluginType = errors.New("cannot parse the spec for an unknown plugin type")

type pluginLoader struct {
	Loader func(root string, params []string) (*Plugin, error)
	Regex  *regexp.Regexp
}

var loaders = []pluginLoader{
	{MakeGitHub, regexp.MustCompile(`github:([a-z0-9\-]+)/([a-z0-9\-]+)`)},
	{MakeDir, regexp.MustCompile(`dir:(.*)`)},
	{MakeFile, regexp.MustCompile(`file:(.*)`)},
	{MakeOhMyZsh, regexp.MustCompile(`oh-my-zsh`)},
	{MakeOhMyZshPlugin, regexp.MustCompile(`oh-my-zsh:plugin:([a-z0-9\-]+)`)},
	{MakeOhMyZshTheme, regexp.MustCompile(`oh-my-zsh:theme:([a-z0-9\-]+)`)},
}

func MakePlugin(root string, spec string) (*Plugin, error) {
	for _, loader := range loaders {
		matches := loader.Regex.FindStringSubmatch(spec)
		if len(matches) > 0 {
			plugin, err := loader.Loader(root, matches[1:])
			if err != nil {
				return nil, errors.Wrap(err, "while loading a plugin")
			}
			return plugin, nil
		}
	}
	return nil, ErrUnknownPluginType
}
