package main

import (
	"context"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/PalPalych7/OtusProjectWork/internal/logger"
	manyarmedbandit "github.com/PalPalych7/OtusProjectWork/internal/manyArmedBandit"
	internalhttp "github.com/PalPalych7/OtusProjectWork/internal/server/http"
	"github.com/PalPalych7/OtusProjectWork/internal/sqlstorage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.toml", "Path to configuration file")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	flag.Parse()
	fmt.Println(flag.Args(), configFile)
	config := NewConfig(configFile)
	fmt.Println("config=", config)
	logg := logger.New(config.Logger.LogFile, config.Logger.Level)
	fmt.Println(config.Logger.Level)
	fmt.Println("logg=", logg)
	logg.Info("Start!")
	myBandid := manyarmedbandit.New(config.Bandit)
	logg.Info("myBandid=", myBandid)

	storage := sqlstorage.New(ctx, config.DB.DBName, config.DB.DBUserName, config.DB.DBPassward, myBandid)
	logg.Info("Get new storage:", storage)
	if err := storage.Connect(); err != nil {
		logg.Fatal(err.Error())
	}
	defer storage.Close()

	server := internalhttp.NewServer(ctx, storage, config.HTTP.Host+":"+config.HTTP.Port, logg)
	defer server.Stop()

	go func() {
		fmt.Println("lets startserver!")
		if err := server.Start(); err != nil {
			logg.Fatal("failed to start http server: " + err.Error())
		}
		<-ctx.Done()
	}()
	<-ctx.Done()
	fmt.Println("finish!")
	logg.Info("finish!")
}
