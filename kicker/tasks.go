package kicker

import (
	"context"
	"kickerbot/db"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	tb "gopkg.in/tucnak/telebot.v3"
)

func TaskKickOldUsers(b tb.Bot) {
	d := db.GetDatabase()
	db.Log("STARTING KICKING TASK", db.EmptyStruct{})
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	now := time.Now().Unix()
	old := now - 60
	filter := bson.D{
		primitive.E{Key: "date_joined", Value: bson.D{{"$lt", old}}},
	}
	users, err := d.GetUsers(ctx, filter)
	if err != nil {
		log.Printf("Error in deleting task: %v", err)
	}
	for _, user := range users {
		chat := tb.Chat{ID: user.ChatId}
		tbUser := tb.User{ID: user.Id}
		member := tb.ChatMember{User: &tbUser}
		b.Ban(&chat, &member)
		d.RemoveUser(ctx, user)
	}
}
