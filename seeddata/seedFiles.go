package seeddata

import (
	"context"
	"fir/models"
	mongoFunctions "fir/mongo"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

// SeedFiles is a list of files to seed the file collection in the database and test file-related functions
var SeedFiles = []interface{}{
	models.FileModel{
		DateCreated:  time.Now().UnixNano(),
		LastModified: time.Now().UnixNano(),
		FilePath:     "/home/aliu/testing.txt",
	},
	models.FileModel{
		DateCreated:  time.Now().UnixNano(),
		LastModified: time.Now().UnixNano(),
		FilePath:     "/home/aliu/lolhello.txt",
	},
}

// InitFiles inserts the seed files specified by seeddata.SeedFiles
func InitFiles(client *mongo.Client, files []interface{}) *mongo.Collection {
	filecollection := mongoFunctions.GetCollection(client, "files")
	insertManyResult, err := filecollection.InsertMany(context.TODO(), files)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertManyResult)
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	return filecollection
}
