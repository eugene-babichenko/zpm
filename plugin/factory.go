package plugin

import (
	"regexp"

	"github.com/pkg/errors"
)

// Returned by `MakePlugin` when it cannot parse a plugin specification string.
var ErrUnknownPluginType = errors.New("cannot parse the spec for an unknown plugin type")

// Specify functions to load a plugin and a regular expression to match
// plugins specification string. The provided function must accept valid regular
// expression matches with the first element (the full string matched) being
// omitted. `root` is the plugin installation directory.
var loaders = []struct {
	Loader func(root string, params map[string]string) (*Plugin, error)
	Regex  *regexp.Regexp
}{
	{MakeGitHub, regexp.MustCompile(`^github:(?P<username>[a-z0-9\-]+)/(?P<repo>[a-z0-9\-]+)(@(?P<version>.+))?$`)},
	{MakeDir, regexp.MustCompile(`^dir:(?P<directory>.*)$`)},
	{MakeFile, regexp.MustCompile(`^file:(?P<filename>.*)$`)},
	{MakeOhMyZsh, regexp.MustCompile(`^oh-my-zsh(@(?P<version>.+))?$`)},
	{MakeOhMyZshPlugin, regexp.MustCompile(`^oh-my-zsh:plugin:(?P<name>[a-z0-9\-]+)$`)},
	{MakeOhMyZshTheme, regexp.MustCompile(`^oh-my-zsh:theme:(?P<name>[a-z0-9\-]+)$`)},
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
			matchesDict := make(map[string]string)
			for idx, match := range matches {
				matchesDict[loader.Regex.SubexpNames()[idx]] = match
			}
			plugin, err := loader.Loader(root, matchesDict)
			if err != nil {
				return nil, errors.Wrap(err, "while loading a plugin")
			}
			return plugin, nil
		}
	}
	return nil, ErrUnknownPluginType
}
