// The detailed documentation on those structs can be found in `README`.
// Explicit YAML annotations are required because go-yaml uses lowercase letters
// by default.
package config

type Logger struct {
	MaxSize    int    `yaml:"MaxSize"`
	MaxAge     int    `yaml:"MaxAge"`
	MaxBackups int    `yaml:"MaxBackups"`
	Level      string `yaml:"Level"`
}

type Config struct {
	Plugins           []string `yaml:"Plugins"`
	Root              string   `yaml:"Root"`
	Logger            Logger   `yaml:"Logger"`
	UpdateCheckPeriod string   `yaml:"UpdateCheckPeriod"`
	LogsPath          string   `yaml:"LogsPath"`
}
