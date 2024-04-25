package config

type Mongo struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
	DB   string `mapstructure:"db"`
}
