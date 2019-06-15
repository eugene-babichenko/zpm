package plugin

type Plugin interface {
	Load() ([]string, error)
	CheckUpdate() (*string, error)
	InstallUpdate() error
}
