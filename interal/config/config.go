package config

import (
	"fmt"
	"sync"

	"github.com/spf13/viper"
	"github.com/zjl410091917/game-rank/interal/utils"
)

type Visitor struct {
	ServerPort int32   `mapstructure:"server_port"`
	Logger     Logger  `mapstructure:"logger"`
	Server     Server  `mapstructure:"server"`
	Mongo      Mongo   `mapstructure:"mongo"`
	Redis      []Redis `mapstructure:"redis"`
}

type Config struct {
	Env  string
	C    Visitor
	Name string
}

var (
	instance *Config
	once     sync.Once
)

func GetInstance() *Config {
	once.Do(func() {
		instance = &Config{}
	})
	return instance
}

func C() *Visitor {
	return &GetInstance().C
}

func (c *Config) OnInit() {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	utils.IfErrPanic(viper.ReadInConfig())
	utils.IfErrPanic(viper.Unmarshal(&c.C))
	fmt.Println(c.C.ServerPort)
}

func (c *Config) OnStart() {
}

func (c *Config) OnStop() {
}
