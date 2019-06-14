package plugin

type Plugin interface {
	Load() ([]string, error)
}
