package main

import (
	"flag"
	"github.com/Sansui233/proxypoolCheck/config"
	"github.com/Sansui233/proxypoolCheck/internal/app"
	"github.com/Sansui233/proxypoolCheck/internal/cron"
	"log"
	"net/http"
)

var configFilePath = ""

func main()  {
	go func() {
		http.ListenAndServe("0.0.0.0:6061", nil)
	}()

	// fetch configuration
	flag.StringVar(&configFilePath, "c", "", "path to config file: config.yaml")
	flag.Parse()
	if configFilePath == "" {
		configFilePath = "config.yaml"
	}
	err := config.Parse(configFilePath)
	if err != nil {
		log.Fatal(err, "\n\"Config file err. Exit\"")
		return
	}

	go app.InitApp()
	go cron.Cron()
	// Run
	app.Run()


}