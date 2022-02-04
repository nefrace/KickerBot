package db

import (
	"context"
	"errors"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func makeIDFilter(user User) bson.D {
	return bson.D{
		primitive.E{Key: "id", Value: user.Id},
		primitive.E{Key: "chat_id", Value: user.ChatId},
	}
}

func (d *DB) UserExists(ctx context.Context, user User) bool {
	filter := makeIDFilter(user)
	return d.EntryExists(ctx, "users", filter)
}

func (d *DB) NewUser(ctx context.Context, user User) error {
	if d.UserExists(ctx, user) {
		return errors.New("user entry already exists")
	}
	err := d.NewEntry(ctx, "users", user)
	if err != nil {
		return err
	}
	log.Printf("New user: %v\n", user)
	return nil
}

func (d *DB) GetUser(ctx context.Context, user User) (User, error) {
	filter := makeIDFilter(user)
	var result User
	err := d.Database.Collection("users").FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return User{}, err
	}
	return result, nil
}

func (d *DB) GetUsers(ctx context.Context, filter bson.D) ([]User, error) {
	cursor, err := d.Database.Collection("users").Find(ctx, filter)
	if err != nil {
		return []User{}, err
	}
	var results []User
	if err = cursor.All(ctx, &results); err != nil {
		return []User{}, err
	}
	return results, nil
}

func (d *DB) RemoveUser(ctx context.Context, user User) error {
	filter := makeIDFilter(user)
	_, err := d.Database.Collection("users").DeleteOne(ctx, filter)
	return err
}
