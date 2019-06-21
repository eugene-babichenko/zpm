package plugin

import (
	"regexp"

	"github.com/pkg/errors"
)

// Returned by `MakePlugin` when it cannot parse a plugin specification string.
var ErrUnknownPluginType = errors.New("cannot parse the spec for an unknown plugin type")

// Specify the function to load a plugin and a regular expression to match
// plugins specification string. The provided function must accept valid regular
// expression matches with the first element (the full string matched) being
// omitted. `root` is the plugin installation directory.
type pluginLoader struct {
	Loader func(root string, params []string) (*Plugin, error)
	Regex  *regexp.Regexp
}

var loaders = []pluginLoader{
	{MakeGitHub, regexp.MustCompile(`^github:([a-z0-9\-]+)/([a-z0-9\-]+)(@(.+))?$`)},
	{MakeDir, regexp.MustCompile(`^dir:(.*)$`)},
	{MakeFile, regexp.MustCompile(`^file:(.*)$`)},
	{MakeOhMyZsh, regexp.MustCompile(`^oh-my-zsh(@(.+))?$`)},
	{MakeOhMyZshPlugin, regexp.MustCompile(`^oh-my-zsh:plugin:([a-z0-9\-]+)$`)},
	{MakeOhMyZshTheme, regexp.MustCompile(`^oh-my-zsh:theme:([a-z0-9\-]+)$`)},
}

// Get the plugin object based on the provided specification. `root` is the
// plugin installation directory. If the provided specification string cannot be
// matched to any known pattern, the `ErrUnknownPluginType` error is returned.
// Accepted plugin types and their specification formats can be found in
// `README.md`.
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
