// The detailed documentation on those structs can be found in `README`.
package config

type Logger struct {
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Level      string
}

type Config struct {
	Plugins           []string
	Root              string
	Logger            Logger
	UpdateCheckPeriod string
	LogsPath          string
}
