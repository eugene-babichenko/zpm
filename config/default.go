package config

var DefaultConfig = Config{
	Plugins: []string{},
	Root:    "",
	Logger: Logger{
		MaxSize:    500,
		MaxAge:     28,
		MaxBackups: 6,
		Level:      "info",
	},
}
