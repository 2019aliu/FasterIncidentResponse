package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// UserModel is the model used for storing users in the database.
// Users are described by multiple fields, including username/password, email, URL, and affiliated group
type UserModel struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Username    string             `json:"username" bson:"username"`
	Password    string             `json:"password" bson:"password"`
	Email       string             `json:"email" bson:"email"`
	URL         string             `json:"url" bson:"url"`
	Groups      []string           `json:"groups" bson:"groups"`
	DateCreated int64              `json:"datecreated,omitempty" bson:"datecreated,omitempty"`
}
