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

type EmptyStruct struct {
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
