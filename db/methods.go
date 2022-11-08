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

func (d *DB) GetChat(ctx context.Context, chatId int64) (Chat, error) {
	filter := bson.D{primitive.E{Key: "id", Value: chatId}}
	var result Chat
	err := d.Database.Collection("chats").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		Log("get chat error", err)
		return Chat{}, err
	}
	return result, nil
}

func (d *DB) NewChat(ctx context.Context, chat Chat) error {
	if d.ChatExists(ctx, chat) {
		return errors.New("chat entry already exists")
	}
	err := d.NewEntry(ctx, "chats", chat)
	if err != nil {
		Log("new chat error", err)
		return err
	}
	log.Printf("New chat added: %s (%d)\n", chat.Title, chat.Id)
	return nil
}

func (d *DB) UpdateChat(ctx context.Context, chat Chat, updates bson.D) error {
	filter := bson.D{primitive.E{Key: "id", Value: chat.Id}}
	result, err := d.Database.Collection("chats").UpdateOne(ctx, filter, updates)
	if err != nil {
		Log("chat update error", err)
	}
	log.Printf("Chat updated: %v", result)
	return err
}