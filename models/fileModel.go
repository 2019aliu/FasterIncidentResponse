package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// FileModel is the model used to represent files in the file collection of the database
type FileModel struct {
	ID           primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DateCreated  int64              `json:"datecreated,omitempty" bson:"datecreated,omitempty"`
	LastModified int64              `json:"lastmodified,omitempty" bson:"lastmodified,omitempty"`
	FilePath     string             `json:"filepath" bson:"filepath"`
}
