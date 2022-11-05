package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	rabbitmq "github.com/PalPalych7/OtusProjectWork/internal/rabbitMQ"
)

type LoggerConf struct {
	LogFile string
	Level   string
}

type DBConf struct {
	DBName     string
	DBUserName string
	DBPassward string
}

type Config struct {
	Logger LoggerConf
	Rabbit rabbitmq.RabbitCFG
	DB     DBConf
}

func NewConfig(configFile string) Config {
	var myConf Config
	_, err := toml.DecodeFile(configFile, &myConf)
	if err != nil {
		fmt.Println("err Decode config File=", err)
	}
	return myConf
}
