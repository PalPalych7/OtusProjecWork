package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	manyarmedbandit "github.com/PalPalych7/OtusProjectWork/internal/manyArmedBandit"
)

type LoggerConf struct {
	LogFile string
	Level   string
}

type HTTPConf struct {
	Host string
	Port string
}

type DBConf struct {
	DBName     string
	DBUserName string
	DBPassward string
}

type Config struct {
	Logger LoggerConf
	HTTP   HTTPConf
	DB     DBConf
	Bandit manyarmedbandit.BanditConfig
}

func NewConfig(configFile string) Config {
	var myConf Config
	_, err := toml.DecodeFile(configFile, &myConf)
	if err != nil {
		fmt.Println("err Decode config File=", err)
	}
	return myConf
}
