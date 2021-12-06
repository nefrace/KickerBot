package main

import (
	"kickerbot/db"
	"kickerbot/kicker"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Print("Error loading .env file")
	}
	token, exists := os.LookupEnv("TOKEN")
	if !exists {
		log.Fatal("no token specified")
	}
	_, dberr := db.Init("mongodb://127.0.0.1:27017")
	if dberr != nil {
		log.Fatal(err)
	}

	Bot := kicker.Kicker{Token: token}
	Bot.Init()
	Bot.AddHandlers(kicker.HandlersV1)
	Bot.Bot.Start()
}
