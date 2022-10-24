package main

import (
	"fmt"

	"github.com/BurntSushi/toml"

	"github.com/PalPalych7/OtusProjectWork/internal/manyArmedBandit"
)

type LoggerConf struct {
	LogFile string
	Level   string
}

type GRPCConf struct {
	Host string
	Port string
}

type DBConf struct {
	DBName     string
	DBUserName string
	DBPassword string
}

type Config struct {
	Logger LoggerConf
	GRPC   GRPCConf
	DB     DBConf
	Bandit manyArmedBandit.BanditConfig
}

func NewConfig(configFile string) Config {
	var myConf Config
	_, err := toml.DecodeFile(configFile, &myConf)
	if err != nil {
		fmt.Println("err Decode config File=", err)
	}
	return myConf
}
