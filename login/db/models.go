package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	// UserID   uint64             `bson:"userid"`
	Username string `bson:"username"`
	Password string `bson:"password"`
}
