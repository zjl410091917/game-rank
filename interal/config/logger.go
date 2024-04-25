package config

type Logger struct {
	Level      string `mapstructure:"level"`
	FileOutput string `mapstructure:"file_output"`
}
