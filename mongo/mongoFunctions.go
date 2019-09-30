package mongo

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var dbname = "test"

// GetClient initializes the MongoDB database, inserting a mongo listener on the default port (27017), and returns a mongo client and
func GetClient() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.NewClient(clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return client
}

// CheckDBConnection ensures the MongoDB client is connected to the database
func CheckDBConnection(client *mongo.Client) {
	err := client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal("Couldn't connect to the database", err)
	} else {
		log.Println("Connected!")
	}
}

// GetCollection retrieves the MongoDB collection for the specified collection.
// In FasterIR, this should be "incidents", "users", "artifacts", or "files"
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	return client.Database(dbname).Collection(collectionName)
}
