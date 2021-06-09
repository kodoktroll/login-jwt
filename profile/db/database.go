package db

import (
	"context"
	"errors"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	collection *mongo.Collection
	Client     *mongo.Client
}

func newDB() *Database {
	db := &Database{}
	return db
}

func SetupDB() (*Database, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()
	mongoUri := os.Getenv("MONGO_URI")
	if len(mongoUri) == 0 {
		mongoUri = "mongodb://localhost:27017"
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	// defer func() {
	// 	if err = client.Disconnect(ctx); err != nil {
	// 		panic(err)
	// 	}
	// }()
	if err != nil {
		log.Print("Error", err.Error())
		return nil, err
	}
	log.Print("Pinging to mongo")
	if err := pingDB(client); err != nil {
		log.Print("Fail to pong")
		pingDB(client)
		// return nil, err
	}
	log.Print("Creating collection")
	collection := client.Database("profile").Collection("profile")
	db := newDB()
	db.collection = collection
	db.Client = client
	return db, nil
}

func (db *Database) InsertProfile(profile Profile) error {
	// defer cancel()
	// log.Print("Inserting profile")
	// var p Profile
	// err := db.collection.FindOne(ctx, bson.M{"userid": profile.UserID}).Decode(&p)
	_, err := db.GetProfile(profile.UserID)
	if err == nil {
		return errors.New("Cannot insert an existing profile")
	}
	_, err = db.collection.InsertOne(context.TODO(), profile)
	if err != nil {
		return err
	}
	log.Print("Profile successfully inserted")
	return nil
}

func (db *Database) GetProfile(userId string) (*Profile, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	// defer cancel()
	var p Profile
	err := db.collection.FindOne(context.TODO(), bson.M{"userid": userId}).Decode(&p)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("No Documents")
		}
		return nil, err
	}
	if len(p.UserID) == 0 {
		return nil, errors.New("Cannot find existing profile")
	}
	return &p, nil
}

func pingDB(client *mongo.Client) error {
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()
	return client.Ping(context.TODO(), readpref.Primary())
}

func (db *Database) CloseDB() error {
	return db.Client.Disconnect(context.TODO())
}
