package main

import (
	"log"
	"os"

	"kickerbot/kicker"
)

func main() {
	token, exists := os.LookupEnv("TOKEN")
	if !exists {
		log.Fatal("no token specified")
	}

	Bot := kicker.Kicker{Token: token}
	Bot.Init()
	Bot.AddHandlers(kicker.HandlersV1)
	Bot.Bot.Start()
}
