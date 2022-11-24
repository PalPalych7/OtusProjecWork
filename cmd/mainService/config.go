package main

import (
	"log"

	"github.com/BurntSushi/toml"
	manyarmedbandit "github.com/PalPalych7/OtusProjectWork/internal/manyArmedBandit"
	"github.com/PalPalych7/OtusProjectWork/internal/sqlstorage"
)

type LoggerConf struct {
	LogFile string
	Level   string
}

type HTTPConf struct {
	Host string
	Port string
}

type Config struct {
	Logger LoggerConf
	HTTP   HTTPConf
	DB     sqlstorage.DBConf
	Bandit manyarmedbandit.BanditConfig
}

func NewConfig(configFile string) Config {
	var myConf Config
	_, err := toml.DecodeFile(configFile, &myConf)
	if err != nil {
		log.Fatal("err Decode config File=", err)
	}
	return myConf
}
