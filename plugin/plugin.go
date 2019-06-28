package plugin

import "github.com/pkg/errors"

// Returned by `Plugin.CheckUpdate` and indicates that a plugin being checked
// have not been installed.
var NotInstalled = errors.New("not installed")

// Returned by `Plugin.CheckUpdate` and indicates that a plugin being checked
// is not upgradable.
var NotUpgradable = errors.New("this plugin is not upgradable")

// Returned by `Plugin.CheckUpdate` and indicates that a plugin being checked
// is up to date.
var UpToDate = errors.New("up to date")

// Check if a plugin is not installed with the error value of
// `Plugin.CheckUpdate`.
func IsNotInstalled(err error) bool {
	return err == NotInstalled
}

// Check if a plugin is cannot be updated with the error value of
// `Plugin.CheckUpdate`.
func IsNotUpgradable(err error) bool {
	return err == NotUpgradable
}

// Check if a plugin is up to date with the error value of
// `Plugin.CheckUpdate`.
func IsUpToDate(err error) bool {
	return err == UpToDate
}

// This is the universal interface for all plugin types loaded by `zpm`. All
// plugins must be used via this interface outside the `plugin` module.
type Plugin interface {
	// Returns `fpath` with functions to be loaded from a plugin and the lines
	// that are required to be executed in order to correctly load a plugin.
	Load() (fpath []string, exec []string, err error)
	// Check for update and return the message describing the update. If no
	// update is available, `nil` is returned instead of an update description.
	// If a plugin is not installed, an error must be set to `NotInstalled` and
	// this can be checked with the `IsNotInstalled` function.
	CheckUpdate() (message *string, err error)
	// Install an update if available or install a plugin if not installed. Must
	// be called after a successful `CheckUpdate` call (that returned an update
	// description or the `NotInstalled` error). Otherwise it may cause a panic
	// or an unexpected error.
	InstallUpdate() error
}
