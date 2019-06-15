package plugin

import "github.com/pkg/errors"

var NotInstalled = errors.New("not installed")

func IsNotInstalled(err error) bool {
	return err == NotInstalled
}

type Plugin interface {
	Load() ([]string, error)
	CheckUpdate() (*string, error)
	InstallUpdate() error
}
