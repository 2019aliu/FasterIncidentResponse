package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// TokenModel is the model used to store tokens in the MongoDB database
type TokenModel struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Token       string             `json:"token" bson:"token"`
	DateCreated int64              `json:"datecreated,omitempty" bson:"datecreated,omitempty"`
}
