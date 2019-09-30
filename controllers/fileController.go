package controllers

import (
	"context"
	"fir/models"
	mongoFunctions "fir/mongo"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

var filecollection = mongoFunctions.GetCollection(mongoFunctions.GetClient(), "files")

// PostFile creates an file and stores it in the FasterIR database
func PostFile(c *gin.Context) {
	var newFile models.FileModel
	reqBody := UnmarshalRequest(c, &newFile)

	if !CheckFile(c, newFile.FilePath) {
		return
	}

	insertResult, err := filecollection.InsertOne(context.TODO(), newFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// Put in initial date of editing
	_, updateerr := filecollection.UpdateOne(context.TODO(), bson.M{"_id": insertResult.InsertedID}, bson.M{"$set": bson.M{"datecreated": time.Now().UnixNano(), "lastmodified": time.Now().UnixNano()}})
	if updateerr != nil {
		log.Fatal(updateerr)
	}

	c.JSON(http.StatusCreated, reqBody)
}

// GetAllFiles retrieves all files stored in the current database
func GetAllFiles(c *gin.Context) {
	var results []*models.FileModel

	cur, err := filecollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.FileModel
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

// GetFile retrieves the file stored in FasterIR for the specified ID
func GetFile(c *gin.Context) {
	id := c.Param("fileID")

	var result models.FileModel
	searchID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": searchID} // search by ID
	err := filecollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)

	c.JSON(http.StatusOK, result)
}

// UpdateFile modifies an file, specified by its ID, for any field provided.
func UpdateFile(c *gin.Context) {
	id := c.Param("fileID")
	searchID, _ := primitive.ObjectIDFromHex(id)

	var updateFields bson.M
	UnmarshalRequest(c, &updateFields)

	if file, fileOK := updateFields["filepath"].(string); fileOK && !CheckFile(c, file) {
		return
	}

	updateBson := bson.M{"$set": updateFields}
	updateResult, err := filecollection.UpdateOne(context.TODO(), bson.M{"_id": searchID}, updateBson, options.Update().SetUpsert(true))
	if err != nil {
		fmt.Println("update error")
		log.Fatal(err)
	}

	_, timeerr := filecollection.UpdateOne(context.TODO(), bson.M{"_id": searchID}, bson.M{"$set": bson.M{"lastmodified": time.Now().UnixNano()}})
	if timeerr != nil {
		log.Fatal(timeerr)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Find the file again
	// I realize this is not the best code ever 'cause boilerplate code but ugha;jas;ldknwdnwjd
	var result models.FileModel
	idfilter := bson.M{"_id": searchID} // search by ID
	finderr := filecollection.FindOne(context.TODO(), idfilter).Decode(&result)
	if finderr != nil {
		log.Fatal(finderr)
	}

	c.JSON(http.StatusCreated, result)
}

// DeleteFile deletes the file corresponding to the specified ID via deleting the matching document in the files collection
func DeleteFile(c *gin.Context) {
	id := c.Param("fileID")
	deletedFileID, _ := primitive.ObjectIDFromHex(id)
	var deletedFile models.FileModel
	err := filecollection.FindOneAndDelete(context.TODO(), bson.M{"_id": deletedFileID}).Decode(&deletedFile)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted Object: %v", deletedFile)

	c.JSON(http.StatusNoContent, deletedFile)
}

// CheckFile verifies if a new or modified filepath should be stored in the database.
// It is used as a helper function for the Create (Post) and Update methods
// Developer's note: a check to differentiate whether the file is already in the database or the file wasn't modified might be useful, but not sure.
func CheckFile(c *gin.Context, filepathname string) bool {
	// Empty filepath check
	if filepathname == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter a filepath.",
		})
		return false
	}
	// Not a valid file path (regex matching)
	if matched, _ := regexp.MatchString(`[^/]+(/[^/]+)*\.[A-Za-z]+`, filepathname); !matched {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Sorry, this is not a valid file path. Make sure your file has a file extension, ensure the path was typed correctly, and please try again!",
		})
		return false
	}
	// Duplicate file check
	cur, _ := filecollection.Find(context.TODO(), bson.M{})
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.FileModel
		cur.Decode(&elem)
		if elem.FilePath == filepathname {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Sorry, either another file with this file path has already been stored or you did not upload a new file. Check your file path and please try again!",
			})
			return false
		}
	}
	// /*Not sure if this works (it should just check whether the file is in the user's file system) */
	// } else if _, fileErr := os.Stat(filepathname); fileErr != nil {
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": "This file does not exist in your file system, please try again!",
	// 	})
	// 	return false
	// }

	return true
}
