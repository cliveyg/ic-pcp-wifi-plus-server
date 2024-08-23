package main

import (
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"os"
)

var kwa []string

func main() {
	kwa = loadKnownWifi()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// logging stuff
	file, err := os.OpenFile(os.Getenv("LOGFILE"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	level, er2 := log.ParseLevel(os.Getenv("LOGLEVEL"))
	if er2 != nil {
		log.Fatal(err)
	}
	log.SetLevel(level)
	log.Info("======================= Starting wifiplus =======================")

	a := App{}
	a.Initialize()

	a.Run(os.Getenv("PORT"))

}
