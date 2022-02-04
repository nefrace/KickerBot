package main

import (
	"kickerbot/captchagen"
	"kickerbot/db"
	"kickerbot/kicker"
	"log"
	"os"
	"time"

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
	Bot.AddHandlers(kicker.HandlersV1)
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.Every(1).Minutes().Do(func() { kicker.TaskKickOldUsers(*Bot.Bot) })
	scheduler.StartAsync()
	Bot.Bot.Start()
}
