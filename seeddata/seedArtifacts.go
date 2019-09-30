package seeddata

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	"fir/models"
	mongoFunctions "fir/mongo"
)

// SeedArtifacts is a list of artifacts used for testing (seeding) functionalities regarding artifacts
var SeedArtifacts = []interface{}{
	models.ArtifactModel{
		Name: "artifact numero uno",
		MArtifact: models.Artifact{
			Description:  "problemo numero uno",
			DateCreated:  time.Now().UnixNano(),
			LastModified: time.Now().UnixNano(),
		},
	},
	models.ArtifactModel{
		Name: "artifact numero dos",
		MArtifact: models.Artifact{
			Description:  "problemo numero dos",
			DateCreated:  time.Now().UnixNano(),
			LastModified: time.Now().UnixNano(),
		},
	},
}

// InitArtifacts inserts the seed artifacts specified by seeddata.SeedArtifacts
func InitArtifacts(client *mongo.Client, artifacts []interface{}) *mongo.Collection {
	artifactcollection := mongoFunctions.GetCollection(client, "artifacts")
	insertManyResult, err := artifactcollection.InsertMany(context.TODO(), artifacts)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertManyResult)
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	return artifactcollection
}
