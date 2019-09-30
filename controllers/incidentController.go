package controllers

import (
	"context"
	"fir/models"
	mongoFunctions "fir/mongo"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var incidentcollection = mongoFunctions.GetCollection(mongoFunctions.GetClient(), "incidents")
var incidentoptionscollection = mongoFunctions.GetCollection(mongoFunctions.GetClient(), "incidentOptions")

// PostIncident creates an incident and stores it in the FasterIR database
func PostIncident(c *gin.Context) {
	var newIncident models.IncidentModel
	reqBody := UnmarshalRequest(c, &newIncident)

	// Check for empty fields
	if newIncident.Detection == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter a detection method.",
		})
		return
	} else if newIncident.Actor == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter an actor.",
		})
		return
	} else if newIncident.Plan == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter a plan.",
		})
		return
	} else if newIncident.Subject == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter a subject.",
		})
		return
	} else if newIncident.Description == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter a description.",
		})
		return
	} else if newIncident.Severity == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter a severity level.",
		})
		return
	} else if newIncident.Status == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter a status.",
		})
		return
	} else if newIncident.Confidentiality == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter a confidentiality level.",
		})
		return
	} else if newIncident.Category == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter a detection method.",
		})
		return
	} else if newIncident.OpenedBy == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter the user that opened this incident.",
		})
		return
	}

	// Reject duplicate subjects
	if !CheckDuplicateSubject(c, newIncident.Subject) {
		return
	}

	// Retrieve the incidentoptions object from the database
	var incidentOptions bson.M
	incidentoptionscollection.FindOne(context.TODO(), bson.M{}).Decode(&incidentOptions)

	// Check whether severity and confidentiality levels are within the scale
	if newIncident.Severity < 1 || newIncident.Severity > incidentOptions["severity"].(int32) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("This severity level is higher than the highest severity level, please enter an integer between 1 and %v", incidentOptions["severity"].(int32)),
		})
		return
	} else if newIncident.Confidentiality < 1 || newIncident.Confidentiality > incidentOptions["confidentiality"].(int32) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("This confidentiality level is higher than the highest confidentiality level, please enter an integer between 1 and %v", incidentOptions["confidentiality"].(int32)),
		})
		return
	}

	// Check whether the fields with options are one of the options
	// Also using deMorgan's law hehe
	if !(CheckIncidentOption(c, "detection", newIncident.Detection) &&
		CheckIncidentOption(c, "actor", newIncident.Actor) &&
		CheckIncidentOption(c, "plan", newIncident.Plan) &&
		CheckIncidentOption(c, "status", newIncident.Status) &&
		CheckIncidentOption(c, "category", newIncident.Category)) {
		return
	}

	// Also need to do the same thing with the concerned business lines
	for _, cbl := range newIncident.ConcernedBusinessLines {
		if !CheckIncidentOption(c, "concernedbusinesslines", cbl) {
			return
		}
	}

	insertResult, err := incidentcollection.InsertOne(context.TODO(), newIncident)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// Put in initial date of editing
	_, updateerr := incidentcollection.UpdateOne(context.TODO(), bson.M{"_id": insertResult.InsertedID}, bson.M{"$set": bson.M{"datecreated": time.Now().UnixNano(), "lastmodified": time.Now().UnixNano()}})
	if updateerr != nil {
		log.Fatal(updateerr)
	}

	c.JSON(http.StatusCreated, reqBody)
}

// GetAllIncidents retrieves all incidents stored in the current database
func GetAllIncidents(c *gin.Context) {
	var results []*models.IncidentModel

	cur, err := incidentcollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.IncidentModel
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}
	// Close the cursor once finished
	cur.Close(context.TODO())
	c.JSON(http.StatusOK, results)
}

// GetIncident retrieves the incident stored in FasterIR for the specified ID
func GetIncident(c *gin.Context) {
	id := c.Param("incidentID")

	var result models.IncidentModel
	searchID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": searchID} // search by ID
	err := incidentcollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)

	c.JSON(http.StatusOK, result)
}

// UpdateIncident modifies an incident, specified by its ID, for any field provided.
// NOTE: Subfields must be updated with dot notation (example: field.subfield: "hello world").
// We are currently trying to fix this problem so proper JSON can be used
func UpdateIncident(c *gin.Context) {
	id := c.Param("incidentID")
	searchID, _ := primitive.ObjectIDFromHex(id)

	var updateFields bson.M
	UnmarshalRequest(c, &updateFields)

	// Reject duplicate subjects
	if subject, subjectOK := updateFields["subject"].(string); subjectOK && !CheckDuplicateSubject(c, subject) {
		return
	}

	// Retrieve the incidentoptions object from the database
	var incidentOptions bson.M
	incidentoptionscollection.FindOne(context.TODO(), bson.M{}).Decode(&incidentOptions)

	// Check whether severity and confidentiality levels are within the scale
	if severity, severityOK := updateFields["severity"].(float64); severityOK && (int32(severity) < 1 || int32(severity) > incidentOptions["severity"].(int32)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("This severity level is not within the range of accepted severity levels, please enter an integer between 1 and %v", incidentOptions["severity"].(int32)),
		})
		return
	} else if confidentiality, confidentialityOK := updateFields["confidentiality"].(float64); confidentialityOK && (int32(confidentiality) < 1 || int32(confidentiality) > incidentOptions["confidentiality"].(int32)) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("This confidentiality level is not within the range of accepted confidentiality levels, please enter an integer between 1 and %v", incidentOptions["confidentiality"].(int32)),
		})
		return
	}

	// Check whether the fields with options are one of the options
	if detection, detectionOK := updateFields["detection"].(string); detectionOK && !CheckIncidentOption(c, "detection", detection) {
		return
	} else if actor, actorOK := updateFields["actor"].(string); actorOK && !CheckIncidentOption(c, "actor", actor) {
		return
	} else if plan, planOK := updateFields["plan"].(string); planOK && !CheckIncidentOption(c, "plan", plan) {
		return
	} else if status, statusOK := updateFields["status"].(string); statusOK && !CheckIncidentOption(c, "status", status) {
		return
	} else if category, categoryOK := updateFields["category"].(string); categoryOK && !CheckIncidentOption(c, "category", category) {
		return
	}

	cblstrings := make([]string, len(updateFields["concernedbusinesslines"].([]interface{})))
	for i, v := range updateFields["concernedbusinesslines"].([]interface{}) {
		cblstrings[i] = fmt.Sprint(v)
	}
	for _, cbl := range cblstrings {
		if !CheckIncidentOption(c, "concernedbusinesslines", cbl) {
			return
		}
	}

	updateBson := bson.M{"$set": updateFields}
	updateResult, err := incidentcollection.UpdateOne(context.TODO(), bson.M{"_id": searchID}, updateBson, options.Update().SetUpsert(true))
	if err != nil {
		fmt.Println("update error")
		log.Fatal(err)
	}

	_, timeerr := incidentcollection.UpdateOne(context.TODO(), bson.M{"_id": searchID}, bson.M{"$set": bson.M{"lastmodified": time.Now().UnixNano()}})
	if timeerr != nil {
		log.Fatal(timeerr)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Find the incident again
	// I realize this is not the best code ever 'cause boilerplate code but ugha;jas;ldknwdnwjd
	var result models.IncidentModel
	idfilter := bson.M{"_id": searchID} // search by ID
	finderr := incidentcollection.FindOne(context.TODO(), idfilter).Decode(&result)
	if finderr != nil {
		log.Fatal(finderr)
	}

	c.JSON(http.StatusCreated, result)
}

// DeleteIncident deletes the incident corresponding to the specified ID via deleting the matching document in the incidents collection
func DeleteIncident(c *gin.Context) {
	id := c.Param("incidentID")
	deletedIncidentID, _ := primitive.ObjectIDFromHex(id)
	var deletedIncident models.IncidentModel
	err := incidentcollection.FindOneAndDelete(context.TODO(), bson.M{"_id": deletedIncidentID}).Decode(&deletedIncident)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted Object: %v", deletedIncident)

	c.JSON(http.StatusNoContent, deletedIncident)
}

// CheckIncidentOption validates whether the new option (newField) for the provided field is one of the current options.
// If not, this method will return false and send a 400 error with an appropriate message.
func CheckIncidentOption(c *gin.Context, field string, newField string) bool {
	var incidentOptions bson.M
	incidentoptionscollection.FindOne(context.TODO(), bson.M{}).Decode(&incidentOptions)

	isOption := false
	for _, option := range incidentOptions[field].(primitive.A) {
		if newField == option.(string) {
			isOption = true
		}
	}
	if !isOption {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("The %v of %v is not a valid %v in the current options. Add this to the %v options or change the entered %v, and try again", field, newField, field, field, field),
		})
		return false
	}
	return true
}

// CheckDuplicateSubject checks if the provided subject is already a subject stored in the incidents collection
func CheckDuplicateSubject(c *gin.Context, newsubject string) bool {
	cur, _ := incidentcollection.Find(context.TODO(), bson.M{})
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.IncidentModel
		cur.Decode(&elem)
		if elem.Subject == newsubject {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Sorry, this subject line has already been taken. Check your file path and please try again!",
			})
			return false
		}
	}
	return true
}

// CheckIntFields validates all integer fields, including severity and confidentiality
// If the provided values are not within their resepective ranges, this method return false
func CheckIntFields(c *gin.Context, newseverity int32, newconfidentiality int32) bool {
	// Retrieve the incidentoptions object from the database
	var incidentOptions bson.M
	incidentoptionscollection.FindOne(context.TODO(), bson.M{}).Decode(&incidentOptions)

	// Check whether severity and confidentiality levels are within the scale
	if newseverity < 1 || newseverity > incidentOptions["severity"].(int32) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("This severity level is higher than the highest severity level, please enter an integer between 1 and %v", incidentOptions["severity"].(int32)),
		})
		return false
	} else if newconfidentiality < 1 || newconfidentiality > incidentOptions["confidentiality"].(int32) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("This confidentiality level is higher than the highest confidentiality level, please enter an integer between 1 and %v", incidentOptions["confidentiality"].(int32)),
		})
		return false
	}
	return true
}
