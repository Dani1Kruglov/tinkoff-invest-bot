package cron

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"os"
	"os/signal"
	"syscall"
	"time"
	"tinkoff-investment-bot/internal/model/tracker"
	"tinkoff-investment-bot/internal/storage"
)

func NewCron(db *gorm.DB, tracker *tracker.Tracker) {
	kldTime, _ := time.LoadLocation("Europe/Kaliningrad")
	scheduler := cron.New(cron.WithLocation(kldTime))

	defer scheduler.Stop()

	_, err := scheduler.AddFunc("* * * * *", func() { //cron at every minute.
		checkShareForPriceDifference(db, tracker)
	})
	if err != nil {
		_ = fmt.Errorf(err.Error())
	}

	go scheduler.Start()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
}

func checkShareForPriceDifference(db *gorm.DB, tracker *tracker.Tracker) {
	shareStorage := storage.NewShareStorage(db)
	shareStorage.GetShares(tracker)
}
