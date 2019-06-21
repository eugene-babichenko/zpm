package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Feature: Load plugins configurations
//   Scenario: load a plugin with Oh My Zsh
//     Given that Oh My Zsh is required by one of the plugins
//     When a user loads the plugins
//     Then Oh My Zsh must be present in the list
//     And be the first
func TestLoadPluginWithOhMyZsh(t *testing.T) {
	var specs = []string{
		"github:zsh-users/zsh-autosuggestions",
		"oh-my-zsh:plugin:pip",
	}

	var expectedNames = []string{
		"oh-my-zsh",
		"github:zsh-users/zsh-autosuggestions",
		"oh-my-zsh:plugin:pip",
	}

	names, _, err := MakePluginsFromSpecs("/root", specs)
	require.Empty(t, err, "cannot parse specs")
	assert.Equal(t, expectedNames, names, "Oh My Zsh must be first in the list")
}

//   Scenario: load a theme with Oh My Zsh
//     Given that Oh My Zsh is required by one of the themes
//     When a user loads the plugins and themes
//     Then Oh My Zsh must be present in the list
//     And be the first
func TestLoadThemeWithOhMyZsh(t *testing.T) {
	var specs = []string{
		"github:zsh-users/zsh-autosuggestions",
		"oh-my-zsh:theme:arrow",
	}

	var expectedNames = []string{
		"oh-my-zsh",
		"github:zsh-users/zsh-autosuggestions",
		"oh-my-zsh:theme:arrow",
	}

	names, _, err := MakePluginsFromSpecs("/root", specs)
	require.Empty(t, err, "cannot parse specs")
	assert.Equal(t, expectedNames, names, "Oh My Zsh must be first in the list")
}

//   Scenario: load  Oh My Zsh
//     Given that Oh My Zsh is in the list
//     When a user loads the plugins and themes
//     Then Oh My Zsh must be present in the list
//     And be the first
func TestLoadOhMyZsh(t *testing.T) {
	var specs = []string{
		"github:zsh-users/zsh-autosuggestions",
		"oh-my-zsh",
	}

	var expectedNames = []string{
		"oh-my-zsh",
		"github:zsh-users/zsh-autosuggestions",
	}

	names, _, err := MakePluginsFromSpecs("/root", specs)
	require.Empty(t, err, "cannot parse specs")
	assert.Equal(t, expectedNames, names, "Oh My Zsh must be first in the list")
}
