package main

import (
	"kickerbot/captchagen"
	"kickerbot/db"
	"kickerbot/kicker"
	"log"
	"os"
	"time"

	"github.com/NicoNex/echotron/v3"
	"github.com/go-co-op/gocron"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	captchagen.Init()
	if err != nil {
		log.Print("Error loading .env file")
	}
	token, exists := os.LookupEnv("TOKEN")
	if !exists {
		log.Fatal("no token specified")
	}
	_, dberr := db.Init(os.Getenv("MONGO_URI"))
	if dberr != nil {
		log.Fatal(err)
	}

	Bot := kicker.Kicker{Token: token}
	Bot.Init()
	scheduler := gocron.NewScheduler(time.UTC)
	tasker := echotron.NewAPI(token)
	scheduler.Every(30).Seconds().Do(func() { kicker.TaskKickOldUsers(&tasker) })
	scheduler.StartAsync()
	Bot.Start()
}
