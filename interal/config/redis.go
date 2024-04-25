package config

type Redis struct {
	Addr    string `mapstructure:"addr"`
	DB      int    `mapstructure:"db"`
	Cluster bool   `mapstructure:"cluster"`
}
