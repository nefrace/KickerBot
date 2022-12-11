package kicker

import (
	"context"
	"kickerbot/db"
	"log"
	"time"

	tb "github.com/NicoNex/echotron/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskBot struct {
	Token string
	tb.API
}

func TaskKickOldUsers(b *tb.API) {
	d := db.GetDatabase()
	log.Print("STARTING KICKING TASK")
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	now := time.Now().Unix()
	old := now - 120
	filter := bson.D{
		primitive.E{Key: "date_joined", Value: bson.D{bson.E{Key: "$lt", Value: old}}},
	}
	users, err := d.GetUsers(ctx, filter)
	if err != nil {
		log.Printf("Error in deleting task: %v", err)
	}
	for _, user := range users {

		_, err := b.BanChatMember(user.ChatId, user.Id, &tb.BanOptions{RevokeMessages: true})
		if err != nil {
			log.Println("User was not banned: ", err)
			continue
		}
		b.DeleteMessage(user.ChatId, user.CaptchaMessage)
		b.DeleteMessage(user.ChatId, user.JoinedMessage)
		d.RemoveUser(ctx, user)
	}
}
