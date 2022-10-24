package main

import (
	"flag"
	"fmt"

	"github.com/PalPalych7/OtusProjectWork/internal/logger"
	"github.com/PalPalych7/OtusProjectWork/internal/manyArmedBandit"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	fmt.Println(flag.Args(), configFile)
	config := NewConfig(configFile)
	fmt.Println("config=", config)
	//	logg := logger.New(config.Logger.LogFile, config.Logger.Level)
	logg := logger.New("", config.Logger.Level)
	fmt.Println(config.Logger.Level)
	fmt.Println("logg=", logg)
	logg.Info("Start!")
	myBandid := manyArmedBandit.New(config.Bandit)
	logg.Info("myBandid=", myBandid, myBandid.GetBannerId(2))
	logg.Info("finish!")
}
