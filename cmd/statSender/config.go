package main

import (
	"fmt"

	"github.com/BurntSushi/toml"
	rabbitmq "github.com/PalPalych7/OtusProjectWork/internal/rabbitMQ"
	"github.com/PalPalych7/OtusProjectWork/internal/sqlstorage"
)

type LoggerConf struct {
	LogFile string
	Level   string
}

type Config struct {
	Logger LoggerConf
	Rabbit rabbitmq.RabbitCFG
	DB     sqlstorage.DBConf
}

func NewConfig(configFile string) Config {
	var myConf Config
	_, err := toml.DecodeFile(configFile, &myConf)
	if err != nil {
		fmt.Println("err Decode config File=", err)
	}
	return myConf
}
