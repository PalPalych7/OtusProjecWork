package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
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
	fmt.Println(flag.Args(), configFile)
	config := NewConfig(configFile)
	fmt.Println("config=", config)
	logg := logger.New(config.Logger.LogFile, config.Logger.Level)
	fmt.Println(config.Logger.Level)
	fmt.Println("logg=", logg)
	logg.Info("Start!")
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()
	fmt.Println("start storage")
	storage := sqlstorage.New(ctx, config.DB, nil)
	logg.Info("Connected to storage:", storage)
	fmt.Println("try connect to storage:", config.DB, storage)
	err := storage.Connect()
	fmt.Println("Connect result:", err)
	if err != nil {
		logg.Fatal(err.Error())
	}
	defer storage.Close()

	fmt.Println("Start Connected to Rabbit:")
	myRQ, err := rabbitmq.CreateQueue(ctx, config.Rabbit)
	fmt.Println("Result Connected to Rabbit:", err)
	if err != nil {
		time.Sleep(time.Minute * 3)
		logg.Fatal(err.Error())
	}
	defer myRQ.Shutdown()

	logg.Info("Connected to Rabit! - ", myRQ)
	go func() {
		for {
			logg.Info("Проснулись.")
			// отправка оповещений
			myStatList, err2 := storage.GetBannerStat()
			countRec := len(myStatList)
			if err2 != nil { //nolint
				logg.Error("ошибка получения статистики-", err2)
			} else if countRec == 0 {
				logg.Info("Данных для отправки не найдено")
			} else {
				logg.Info("Найдено ", countRec, "записей для отправки")
				myMess, errMarsh := json.Marshal(myStatList)
				if errMarsh != nil {
					logg.Error("ошибка json.Marshal", errMarsh)
				}
				if erSemdMess := myRQ.SendMess(myMess); erSemdMess != nil {
					logg.Error("ошибка отправки сообщения-", errMarsh)
				} else {
					logg.Info("сообщение успешно отпралвено")
				}
				myStatID := myStatList[countRec-1].ID
				logg.Info("max_stat_id=", myStatID)
				if errChID := storage.ChangeSendStatID(myStatID); errChID != nil {
					logg.Error("ошибка обновления max ID отправки -", errMarsh)
				}
			}
			time.Sleep(time.Minute * 5)
		}
	}()
	<-ctx.Done()
	logg.Info("Finish")
}
