package db

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (d *DB) GetEntry(ctx context.Context, collectionName string, filter interface{}) (interface{}, error) {
	collection := d.Database.Collection(collectionName)
	var result bson.D
	err := collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return EmptyStruct{}, err
	}
	return result, nil
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
	_, err := d.GetEntry(ctx, collectionName, filter)
	if err != nil {
		log.Printf("EntryExists error: %v", err)
	}
	return err != mongo.ErrNoDocuments
}
