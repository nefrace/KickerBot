package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	Initialized bool
	Ctx         context.Context
	URI         string
	Database    *mongo.Database
	Client      *mongo.Client
}

var database DB = DB{}

func Init(URI string) (DB, error) {
	if (DB{}) != database {
		return DB{}, errors.New("database already exists")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URI).SetAuth(options.Credential{Username: os.Getenv("DB_USER"), Password: os.Getenv("DB_PASSWORD")}))
	if err != nil {
		return DB{}, err
	}
	db := client.Database("godotkicker")
	database = DB{
		URI:      URI,
		Database: db,
		Client:   client,
	}
	Log("startup", nil)
	return database, nil
}

func Log(event string, what interface{}) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	t := time.Now()
	content := fmt.Sprintf("%v", what)
	log.Printf("%v :: EVENT: %s\n%v", t, event, what)
	_, err := database.Database.Collection("logs").InsertOne(ctx, bson.D{
		primitive.E{Key: "time", Value: t},
		primitive.E{Key: "event", Value: event},
		primitive.E{Key: "what", Value: content},
	})
	if err != nil {
		log.Print(err)
	}

}

func GetDatabase() DB {
	return database
}

func (d *DB) Stop() {
	if err := d.Client.Disconnect(d.Ctx); err != nil {
		panic(err)
	}
}

func (d *DB) NewEntry(ctx context.Context, collectionName string, entry interface{}) error {
	collection := d.Database.Collection(collectionName)
	res, err := collection.InsertOne(ctx, entry)
	if err != nil {
		return err
	}
	log.Printf("New entry: %v\nDB ID: %v\n", entry, res.InsertedID)
	return nil
}

func (d *DB) EntryExists(ctx context.Context, collectionName string, filter interface{}) bool {
	collection := d.Database.Collection(collectionName)
	var result bson.D
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		log.Printf("EntryExists error: %v", err)
	}
	return err != mongo.ErrNoDocuments
}

func (d *DB) ChatExists(ctx context.Context, chat Chat) bool {
	filter := bson.D{primitive.E{Key: "id", Value: chat.Id}}
	return d.EntryExists(ctx, "chats", filter)
}

func (d *DB) UserExists(ctx context.Context, user User) bool {
	filter := bson.D{
		primitive.E{Key: "id", Value: user.Id},
		primitive.E{Key: "chat_id", Value: user.ChatId},
	}
	return d.EntryExists(ctx, "users", filter)
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
