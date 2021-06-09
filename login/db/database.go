package db

import (
	"context"
	"errors"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Database struct {
	collection *mongo.Collection
	client     *mongo.Client
}

func newDB() *Database {
	db := &Database{}
	return db
}

func (d *Database) GetUser(username string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	var u User
	err := d.collection.FindOne(ctx, bson.M{"username": username}).Decode(&u)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("No Document")
		}
		return nil, err
	}
	return &u, nil
}

func (d *Database) InsertUser(user User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	log.Print("Inserting User")
	var u User
	err := d.collection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&u)
	if err == nil {
		return errors.New("Username has been used")
	}
	_, err = d.collection.InsertOne(ctx, user)
	if err != nil {
		log.Print("Failed to insert user, ", err.Error())
		return err
	}
	log.Print("User successfully inserted")
	return nil
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
	if err := PingDB(client); err != nil {
		log.Print("Fail to pong")
		PingDB(client)
		// return nil, err
	}
	log.Print("Creating collection")
	collection := client.Database("user").Collection("user")
	db := newDB()
	db.collection = collection
	db.client = client
	return db, nil
}

func PingDB(client *mongo.Client) error {
	// ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	// defer cancel()
	return client.Ping(context.TODO(), readpref.Primary())
}

func (d *Database) CloseDB() error {
	return d.client.Disconnect(context.TODO())
}
