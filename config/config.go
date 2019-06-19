package config

type Logger struct {
	MaxSize    int
	MaxAge     int
	MaxBackups int
	Level      string
}

type Config struct {
	Plugins []string
	Root    string
	Logger  Logger
}
