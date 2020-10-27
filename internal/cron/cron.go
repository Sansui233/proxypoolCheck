package cron

import (
	"github.com/Sansui233/proxypoolCheck/config"
	"github.com/Sansui233/proxypoolCheck/internal/app"
	"github.com/jasonlvhit/gocron"
	"log"
	"runtime"
	"time"
)

func Cron() {
	_ = gocron.Every(15).Minutes().Do(appTask)
	<-gocron.Start()
}

func appTask() {
	config.Parse("")
	err := app.InitApp()
	if err != nil { // for wake up heroku
		log.Println("Init app err: ", err, "\n Try in 2 minute")
		time.Sleep(time.Minute*2)
		app.InitApp()
	}
	runtime.GC()
}
