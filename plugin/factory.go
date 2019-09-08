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
	Loader func(string, map[string]string) (*Plugin, error)
	Regex  *regexp.Regexp
}{
	{MakeGitHub, regexp.MustCompile(`^github\.com/(?P<username>[a-z0-9\-]+)/(?P<repo>[a-z0-9\-]+)(@(?P<version>.+))?$`)},
	{MakeDir, regexp.MustCompile(`^dir://(?P<directory>.*)$`)},
}

var ohMyZshRegex = regexp.MustCompile(`^oh-my-zsh(@(?P<version>.+))?$`)

var ohMyZshLoaders = []struct {
	Loader func(OhMyZsh, map[string]string) (*Plugin, error)
	Regex  *regexp.Regexp
}{
	{MakeOhMyZshPlugin, regexp.MustCompile(`^oh-my-zsh/plugin/(?P<name>[a-z0-9\-]+)$`)},
	{MakeOhMyZshTheme, regexp.MustCompile(`^oh-my-zsh/theme/(?P<name>[a-z0-9\-]+)$`)},
}

type Factory struct {
	Root        string
	ohMyZsh     *OhMyZsh
	ohMyZshSpec string
}

// Get the plugin object based on the provided specification. `root` is the
// plugin installation directory. If the provided specification string cannot be
// matched to any known pattern, the `ErrUnknownPluginType` error is returned.
// Accepted plugin types and their specification formats can be found in
// `README.md`.
func (f *Factory) MakePlugin(spec string) (*Plugin, bool, error) {
	for _, loader := range loaders {
		matches := loader.Regex.FindStringSubmatch(spec)
		if len(matches) == 0 {
			continue
		}
		matchesDict := make(map[string]string)
		for idx, match := range matches {
			matchesDict[loader.Regex.SubexpNames()[idx]] = match
		}
		plugin, err := loader.Loader(f.Root, matchesDict)
		if err != nil {
			return nil, false, errors.Wrap(err, "while loading a plugin")
		}
		return plugin, false, nil
	}

	for _, loader := range ohMyZshLoaders {
		matches := loader.Regex.FindStringSubmatch(spec)
		if len(matches) == 0 {
			continue
		}
		matchesDict := make(map[string]string)
		for idx, match := range matches {
			matchesDict[loader.Regex.SubexpNames()[idx]] = match
		}
		if f.ohMyZsh == nil {
			ohMyZsh := MakeOhMyZsh(f.Root, make(map[string]string))
			f.ohMyZshSpec = "oh-my-zsh"
			f.ohMyZsh = &ohMyZsh
		}
		plugin, err := loader.Loader(*f.ohMyZsh, matchesDict)
		if err != nil {
			return nil, false, errors.Wrap(err, "while loading a plugin")
		}
		return plugin, false, nil
	}

	matches := ohMyZshRegex.FindStringSubmatch(spec)
	if len(matches) == 0 {
		return nil, false, ErrUnknownPluginType
	}

	if f.ohMyZsh != nil {
		plugin := Plugin(f.ohMyZsh)
		return &plugin, true, nil
	}
	matchesDict := make(map[string]string)
	for idx, match := range matches {
		matchesDict[ohMyZshRegex.SubexpNames()[idx]] = match
	}
	ohMyZsh := MakeOhMyZsh(f.Root, matchesDict)
	plugin := Plugin(f.ohMyZsh)
	f.ohMyZshSpec = spec
	f.ohMyZsh = &ohMyZsh
	return &plugin, true, nil
}

func (f *Factory) Dependencies() ([]Plugin, []string) {
	if f.ohMyZsh != nil {
		return []Plugin{Plugin(f.ohMyZsh)}, []string{f.ohMyZshSpec}
	}
	return nil, nil
}
