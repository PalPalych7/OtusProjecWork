package main

import (
	"context"
	"encoding/json"
	"flag"
	"os/signal"
	"syscall"
	"time"

	"github.com/PalPalych7/OtusProjectWork/internal/logger"
	rabbitmq "github.com/PalPalych7/OtusProjectWork/internal/rabbitMQ"
	"github.com/PalPalych7/OtusProjectWork/internal/sqlstorage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/statSenderConfig.toml", "Path to configuration file")
}

func main() {
	flag.Parse()
	config := NewConfig(configFile)
	logg := logger.New(config.Logger.LogFile, config.Logger.Level)
	logg.Info("Start!")
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	storage := sqlstorage.New(ctx, config.DB, nil)
	logg.Info("Connected to storage:", storage)
	if err := storage.Connect(); err != nil {
		logg.Fatal(err.Error())
	}
	defer storage.Close()
	myRQ, err := rabbitmq.CreateQueue(ctx, config.Rabbit)
	if err != nil {
		time.Sleep(time.Minute * 1)
		logg.Fatal(err.Error())
	}
	defer myRQ.Shutdown()

	logg.Info("Connected to Rabit! - ", myRQ)
	go func() {
		for {
			logg.Info("I not sleep :).")
			// отправка оповещений
			myStatList, err2 := storage.GetBannerStat()
			countRec := len(myStatList)
			switch {
			case err2 != nil:
				logg.Error("Error in GetBannerStat", err2)
			case countRec == 0:
				logg.Info("Nothing found for sending")
			default:
				logg.Info("Found ", countRec, "record for sending")
				myMess, errMarsh := json.Marshal(myStatList)
				if errMarsh != nil {
					logg.Error("json.Marshal error", errMarsh)
				}
				if erSemdMess := myRQ.SendMess(myMess); erSemdMess != nil {
					logg.Error("Send mesage error", errMarsh)
				} else {
					logg.Info("message was succcessful send")
				}
				myStatID := myStatList[countRec-1].ID
				logg.Info("max_stat_id=", myStatID)
				if errChID := storage.ChangeSendStatID(myStatID); errChID != nil {
					logg.Error("error in update max send ID -", errMarsh)
				}
			}
			time.Sleep(time.Minute * time.Duration(config.Rabbit.SleepMinutes))
		}
	}()
	<-ctx.Done()
	logg.Info("Finish")
}
