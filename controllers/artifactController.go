/*
Package controllers contains the CRUD operations for interacting with the database (currently MongoDB).
Developer note: the code can be cleaner by removing repetitive code
*/
package controllers

import (
	"context"
	"fir/models"
	mongoFunctions "fir/mongo"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"gopkg.in/mgo.v2/bson"
)

var artifactcollection = mongoFunctions.GetCollection(mongoFunctions.GetClient(), "artifacts")

// PostArtifact creates an artifact and stores it in the FasterIR database
func PostArtifact(c *gin.Context) {
	var newArtifact models.ArtifactModel
	reqBody := UnmarshalRequest(c, &newArtifact)

	if !CheckArtifactName(c, newArtifact.Name) {
		return
	}

	insertResult, err := artifactcollection.InsertOne(context.TODO(), newArtifact)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// Put in initial date of editing
	_, updateerr := artifactcollection.UpdateOne(context.TODO(), bson.M{"_id": insertResult.InsertedID}, bson.M{"$set": bson.M{"mArtifact.datecreated": time.Now().UnixNano(), "mArtifact.lastmodified": time.Now().UnixNano()}})
	if updateerr != nil {
		log.Fatal(updateerr)
	}

	c.JSON(http.StatusCreated, reqBody)
}

// GetAllArtifacts retrieves all artifacts stored in the FasterIR database
func GetAllArtifacts(c *gin.Context) {
	var results []*models.ArtifactModel

	cur, err := artifactcollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.ArtifactModel
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

// GetArtifact retrieves the artifact stored in FasterIR for the specified ID
func GetArtifact(c *gin.Context) {
	id := c.Param("artifactID")

	var result models.ArtifactModel
	searchID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": searchID} // search by ID
	err := artifactcollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)

	c.JSON(http.StatusOK, result)
}

// UpdateArtifact modifies an artifact, specified by its ID, for any field provided.
// NOTE: Subfields must be updated with dot notation (example: field.subfield: "hello world").
// We are currently trying to fix this problem so proper JSON can be used
func UpdateArtifact(c *gin.Context) {
	id := c.Param("artifactID")
	searchID, _ := primitive.ObjectIDFromHex(id)

	var updateFields bson.M
	UnmarshalRequest(c, &updateFields)

	if artifactName, artifactNameOK := updateFields["name"].(string); artifactNameOK && !CheckArtifactName(c, artifactName) {
		return
	}

	updateBson := bson.M{"$set": updateFields}
	updateResult, err := artifactcollection.UpdateOne(context.TODO(), bson.M{"_id": searchID}, updateBson, options.Update().SetUpsert(true))
	if err != nil {
		fmt.Println("update error")
		log.Fatal(err)
	}

	_, timeerr := artifactcollection.UpdateOne(context.TODO(), bson.M{"_id": searchID}, bson.M{"$set": bson.M{"mArtifact.lastmodified": time.Now().UnixNano()}})
	if timeerr != nil {
		log.Fatal(timeerr)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Find the artifact again
	// I realize this is not the best code ever 'cause boilerplate code but ugha;jas;ldknwdnwjd
	var result models.ArtifactModel
	idfilter := bson.M{"_id": searchID} // search by ID
	finderr := artifactcollection.FindOne(context.TODO(), idfilter).Decode(&result)
	if finderr != nil {
		log.Fatal(finderr)
	}

	c.JSON(http.StatusCreated, result)
}

// DeleteArtifact deletes the artifact corresponding to the specified ID via deleting the matching document in the artifacts collection
func DeleteArtifact(c *gin.Context) {
	id := c.Param("artifactID")
	deletedArtifactID, _ := primitive.ObjectIDFromHex(id)
	var deletedArtifact models.ArtifactModel
	err := artifactcollection.FindOneAndDelete(context.TODO(), bson.M{"_id": deletedArtifactID}).Decode(&deletedArtifact)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted artifact: %v", deletedArtifact)

	c.JSON(http.StatusNoContent, deletedArtifact)
}

// CheckArtifactName validates whether an artifact's name is non-empty and is unique in the database
// If either one of those conditions are not satisfied, this method returns false; otherwise, it returns true
// This is used as a helper function for creating and updating artifacts
func CheckArtifactName(c *gin.Context, artifactName string) bool {
	if artifactName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter a name for the artifact",
		})
		return false
	}
	// Duplicate file check
	cur, _ := artifactcollection.Find(context.TODO(), bson.M{})
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.ArtifactModel
		cur.Decode(&elem)
		if elem.Name == artifactName {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Sorry, either another artifact with this name has been stored or the artifact name hasn't been updated. Check the artifact name and please try again!",
			})
			return false
		}
	}
	return true
}
