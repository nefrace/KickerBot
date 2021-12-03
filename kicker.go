package main

import (
	"log"
	"os"

	"godotkicker/bot"
)

func main() {
	token, exists := os.LookupEnv("TOKEN")
	if !exists {
		log.Fatal("no token specified")
	}

	Bot := bot.Bot{Token: token}
	Bot.Init()
	Bot.AddHandlers(bot.HandlersV1)

	log.Print("successfuly launched")
	Bot.Bot.Start()
}
