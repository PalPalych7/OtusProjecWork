package main

import (
	"context"
	"flag"
	"os/signal"
	"syscall"
	"time"

	"github.com/PalPalych7/OtusProjectWork/internal/logger"
	manyarmedbandit "github.com/PalPalych7/OtusProjectWork/internal/manyArmedBandit"
	internalhttp "github.com/PalPalych7/OtusProjectWork/internal/server/http"
	"github.com/PalPalych7/OtusProjectWork/internal/sqlstorage"
)

var configFile string

type serverInt interface {
	Serve() error
	Stop() error
}

var server serverInt

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.toml", "Path to configuration file")
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	flag.Parse()
	config := NewConfig(configFile)
	logg := logger.New(config.Logger.LogFile, config.Logger.Level)
	logg.Info("Start!")
	myBandid := manyarmedbandit.New(config.Bandit)

	storage := sqlstorage.New(ctx, config.DB, myBandid)
	if err := storage.Connect(); err != nil {
		time.Sleep(time.Minute * 1)
		logg.Fatal(err.Error())
	}
	defer storage.Close()

	server = internalhttp.NewServer(ctx, storage, ":"+config.HTTP.Port, logg)
	defer server.Stop()

	go func() {
		if err := server.Start(); err != nil {
			logg.Fatal("failed to start http server: " + err.Error())
		} else {
			logg.Info("Server was started")
		}
		<-ctx.Done()
	}()
	<-ctx.Done()
	logg.Info("Finish!")
}
