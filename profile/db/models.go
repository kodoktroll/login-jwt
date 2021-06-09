package db

import "go.mongodb.org/mongo-driver/bson/primitive"

type Profile struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	UserID      string             `bson:"userid"`
	Username    string             `bson:"username"`
	FirstName   string             `bson:"firstname"`
	LastName    string             `bson:"lastname"`
	PictureLink string             `bson:"pictlink"`
}
