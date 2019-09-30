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

// SeedUsers is the list of users used to seed the users collection in the database and test out user functionalities
var SeedUsers = []interface{}{
	models.UserModel{
		Username:    "2019aliu",
		Password:    "alexliupassword",
		Email:       "alex.liu@fluencysecurity.com",
		URL:         "https://user.tjhsst.edu/2019aliu/",
		Groups:      []string{"interns", "unit_testing", "fasterIR"},
		DateCreated: time.Now().UnixNano(),
	},
	models.UserModel{
		Username:    "2020lluo",
		Password:    "lukeluopassword",
		Email:       "luke.luo@fluencysecurity.com",
		URL:         "https://www.google.com/",
		Groups:      []string{"interns", "fasterIR"},
		DateCreated: time.Now().UnixNano(),
	},
}

// InitUsers inserts the seed users specified by seeddata.SeedUsers
func InitUsers(client *mongo.Client, users []interface{}) *mongo.Collection {
	usercollection := mongoFunctions.GetCollection(client, "users")
	insertManyResult, err := usercollection.InsertMany(context.TODO(), users)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertManyResult)
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	return usercollection
}
