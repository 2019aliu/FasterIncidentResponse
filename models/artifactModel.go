/*
Package models contains the schemas to be used whe storing the different types of objects.
Golang's driver allows for the models to be defined as structs.
*/
package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	// Artifact is anything that is potential evidence
	Artifact struct {
		Description  string `json:"description" bson:"description"`
		DateCreated  int64  `json:"datecreated,omitempty" bson:"datecreated,omitempty"`
		LastModified int64  `json:"lastmodified,omitempty" bson:"lastmodified,omitempty"`
	}

	// ArtifactModel is the model used for all artifacts, contains name and the artifact itself
	ArtifactModel struct {
		ID        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
		Name      string             `json:"name" bson:"name"`
		MArtifact Artifact           `json:"mArtifact" bson:"mArtifact"`
	}
)
