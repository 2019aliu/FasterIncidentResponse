package seeddata

import (
	"context"
	"fir/models"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	mongoFunctions "fir/mongo"
)

// SeedIncidents is a list of incidents that is used to populating the incidents collection and testing functionalities regarding incidents
var SeedIncidents = []interface{}{
	models.IncidentModel{
		Detection:              "cert",
		Actor:                  "cert",
		Plan:                   "a",
		FileSet:                []string{"testing.txt"},
		DateCreated:            time.Now().UnixNano(),
		LastModified:           time.Now().UnixNano(),
		IsStarred:              false,
		Subject:                "testingtesting",
		Description:            "testingtestinghuhhhhhhh",
		Severity:               1,
		IsIncident:             true,
		IsMajor:                true,
		Status:                 "open",
		Confidentiality:        1,
		Category:               "phishing",
		OpenedBy:               "admin",
		ConcernedBusinessLines: []string{"cert"},
	},
	models.IncidentModel{
		Detection:              "external",
		Actor:                  "entity",
		Plan:                   "a",
		FileSet:                []string{"myfriendlyvirus.txt"},
		DateCreated:            time.Now().UnixNano(),
		LastModified:           time.Now().UnixNano(),
		IsStarred:              true,
		Subject:                "yo hablo espanol did you know that",
		Description:            "duolingo makes you an expert at any language lol",
		Severity:               2,
		IsIncident:             true,
		IsMajor:                false,
		Status:                 "open",
		Confidentiality:        2,
		Category:               "scam_web",
		OpenedBy:               "dev",
		ConcernedBusinessLines: []string{"demobl1"},
	},
}

// InitIncidents inserts the seed incidents specified by seeddata.SeedIncidents
func InitIncidents(client *mongo.Client, incidents []interface{}) *mongo.Collection {
	incidentcollection := mongoFunctions.GetCollection(client, "incidents")
	insertManyResult, err := incidentcollection.InsertMany(context.TODO(), incidents)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(insertManyResult)
	fmt.Println("Inserted multiple documents: ", insertManyResult.InsertedIDs)

	return incidentcollection
}

// IncidentOptions are the different options that are available for each category
var IncidentOptions = bson.M{"$set": bson.M{
	"detection": []string{
		"cert",
		"external",
		"bl",
		"soc",
	},
	"actor": []string{
		"cert",
		"entity",
	},
	"plan": []string{
		"a",
		"b",
		"c",
		"1",
		"2",
		"5",
		"6",
	},
	"severity": 4, // Maximum severity
	"status": []string{
		"open",
		"closed",
		"blocked",
	},
	"confidentiality": 4, // Maximum confidentiality
	"category": []string{
		"phishing",
		"scam_web",
		"malware",
		"dataleak",
		"cybersquatting",
		"stolen_data",
		"scam_msg",
		"unavailability",
		"is_integrity",
		"fraud",
		"compromise",
		"reputation",
		"vulnerability",
		"spam",
		"social_eng",
		"consulting",
		"threatIntel",
		"insider",
		"blackmail",
		"dos",
		"scam_tel",
		"scam_social",
		"security_assessment",
	},
	"concernedbusinesslines": []string{
		"cert",
		"demobl1",
		"demobl1_subbl1",
		"demobl1_subbl2",
		"demobl2",
	},
}}

// InitIncidentOptions initializes the options for the incidents by upserting the options into a MongoDB collection
func InitIncidentOptions(client *mongo.Client) *mongo.Collection {
	incidentoptionscollection := mongoFunctions.GetCollection(client, "incidentOptions")
	_, err := incidentoptionscollection.UpdateOne(context.TODO(), bson.M{}, IncidentOptions, options.Update().SetUpsert(true))
	if err != nil {
		log.Fatal(err)
	}

	return incidentoptionscollection
}
