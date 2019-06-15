package plugin

import "github.com/pkg/errors"

var NotInstalledError = errors.New("not installed")

func IsNotInstalled(err error) bool {
	return err == NotInstalledError
}

type Plugin interface {
	Load() ([]string, error)
	CheckUpdate() (*string, error)
	InstallUpdate() error
}
