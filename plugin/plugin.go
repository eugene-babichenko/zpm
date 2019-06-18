package plugin

import "github.com/pkg/errors"

var NotInstalled = errors.New("not installed")

func IsNotInstalled(err error) bool {
	return err == NotInstalled
}

type Plugin interface {
	Load() (fpath []string, exec []string, err error)
	CheckUpdate() (message *string, err error)
	InstallUpdate() error
	GetPath() string
}
