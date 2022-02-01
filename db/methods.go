package db

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (d *DB) ChatExists(ctx context.Context, chat Chat) bool {
	filter := bson.D{primitive.E{Key: "id", Value: chat.Id}}
	return d.EntryExists(ctx, "chats", filter)
}

func (d *DB) NewChat(ctx context.Context, chat Chat) error {
	if d.ChatExists(ctx, chat) {
		return errors.New("chat entry already exists")
	}
	err := d.NewEntry(ctx, "chats", chat)
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.Printf("New chat added: %s (%d)\n", chat.Title, chat.Id)
	return nil
}
