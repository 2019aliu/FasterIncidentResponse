package controllers

// There is a LOT of repetition in the fir code. Don't worry, this is so we can test the functionalities more easily

import (
	"context"
	"fir/models"
	mongoFunctions "fir/mongo"
	redisFunctions "fir/redis"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

var usercollection = mongoFunctions.GetCollection(mongoFunctions.GetClient(), "users")
var saltRounds = 20

// PostUser creates an user and stores it in the FasterIR database
func PostUser(c *gin.Context) {
	var newUser models.UserModel
	reqBody := UnmarshalRequest(c, &newUser)

	// Check for missing fields
	if newUser.Username == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter a username.",
		})
		return
	} else if newUser.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter a password.",
		})
		return
	} else if newUser.Email == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter a email.",
		})
		return
	} else if newUser.URL == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Please enter a URL.",
		})
		return
	}

	// Validate the user
	if !ValidateUser(c, newUser.Username, newUser.Email, newUser.URL) {
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(newUser.Password), saltRounds)
	newUser.Password = string(hashedPassword)

	insertResult, err := usercollection.InsertOne(context.TODO(), newUser)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Inserted a single document: ", insertResult.InsertedID)

	// Put in initial date of editing
	_, updateerr := usercollection.UpdateOne(context.TODO(), bson.M{"_id": insertResult.InsertedID}, bson.M{"$set": bson.M{"datecreated": time.Now().UnixNano()}})
	if updateerr != nil {
		log.Fatal(updateerr)
	}

	c.JSON(http.StatusCreated, reqBody)
}

// GetAllUsers retrieves all users stored in the FasterIR database
func GetAllUsers(c *gin.Context) {
	var req *http.Request = c.Request
	cookie, err := req.Cookie("sessiontoken")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			c.String(http.StatusUnauthorized, "The cookie was not set")
			return
		}
		// For any other type of error, return a bad request status
		c.String(http.StatusBadRequest, "Something happened in between :( please try again")
		return
	}
	sessionToken := cookie

	var cache = redisFunctions.InitCache()
	// We then get the name of the user from our cache, where we set the session token
	response, err := cache.Do("GET", sessionToken)
	if err != nil {
		// If there is an error fetching from cache, return an internal server error status
		c.String(http.StatusInternalServerError, "Server error with get tokens")
		return
	}
	if response == nil {
		// If the session token is not present in cache, return an unauthorized error
		c.String(http.StatusUnauthorized, "There is no session token in the redis cache, try again?")
		return
	}

	var results []*models.UserModel

	cur, err := usercollection.Find(context.TODO(), bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	// Finding multiple documents returns a cursor
	// Iterating through the cursor allows us to decode documents one at a time
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.UserModel
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

// GetUser retrieves the user stored in FasterIR for the specified ID
func GetUser(c *gin.Context) {
	id := c.Param("userID")

	var result models.UserModel
	searchID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": searchID} // search by ID
	err := usercollection.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found a single document: %+v\n", result)

	c.JSON(http.StatusOK, result)
}

// UpdateUser modifies an user, specified by its ID, for any field provided.
func UpdateUser(c *gin.Context) {
	id := c.Param("userID")
	searchID, _ := primitive.ObjectIDFromHex(id)

	var updateFields bson.M
	UnmarshalRequest(c, &updateFields)

	// Field validations
	if username, usernameOK := updateFields["username"].(string); usernameOK && !ValidateUsername(c, username) {
		return
	} else if email, emailOK := updateFields["email"].(string); emailOK && !ValidateEmail(c, email) {
		return
	} else if url, urlOK := updateFields["url"].(string); urlOK && !ValidateURL(c, url) {
		return
	}

	updateBson := bson.M{"$set": updateFields}
	updateResult, err := usercollection.UpdateOne(context.TODO(), bson.M{"_id": searchID}, updateBson, options.Update().SetUpsert(true))
	if err != nil {
		fmt.Println("update error")
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", updateResult.MatchedCount, updateResult.ModifiedCount)

	// Find the user again
	// I realize this is not the best code ever 'cause boilerplate code but ugha;jas;ldknwdnwjd
	var result models.UserModel
	idfilter := bson.M{"_id": searchID} // search by ID
	finderr := usercollection.FindOne(context.TODO(), idfilter).Decode(&result)
	if finderr != nil {
		log.Fatal(finderr)
	}

	c.JSON(http.StatusCreated, result)
}

// DeleteUser deletes the user corresponding to the specified ID via deleting the matching document in the users collection
func DeleteUser(c *gin.Context) {
	id := c.Param("userID")
	deletedUserID, _ := primitive.ObjectIDFromHex(id)
	var deletedUser models.UserModel
	err := usercollection.FindOneAndDelete(context.TODO(), bson.M{"_id": deletedUserID}).Decode(&deletedUser)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Deleted Object: %v", deletedUser)

	c.JSON(http.StatusNoContent, deletedUser)
}

// ValidateUser validates that the information provide conforms to FIR guidelines.
// It performs checks on username, email address, and URL, making sure all are valid.
// If any one of the fields are not validated, this function will return false; otherwise, it will return true.
// It is used as a helper function for creating and updating users
func ValidateUser(c *gin.Context, username string, email string, URL string) bool {
	return ValidateUsername(c, username) && ValidateEmail(c, email) && ValidateURL(c, URL)
}

// ValidateUsername ensures the username:
// 1) does not contain pluses (+), commas (,), brackets (<>), equals (=), or ampersands (&).
// 2) contains between 4 to 32 characters, inclusive. " +
// 3) starts with an alphanumeric character (A-Z, a-z, 0-9, _)
// 4) is not a duplicate username in the database (users collection)
// If the username does not uphold any one of the four rules, this method returns false; otherwise it returns true.
func ValidateUsername(c *gin.Context, username string) bool {
	// Invalid username check
	if validUsername, _ := regexp.MatchString(`^\w(?:[^+,<>=&]){3,32}$`, username); !validUsername {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "This username is invalid. Check to ensure your username: " +
				"1) does not contain pluses (+), commas (,), brackets (<>), equals (=), or ampersands (&). " +
				"2) contains between 4 to 32 characters, inclusive. " +
				"3) starts with an alphanumeric character (A-Z, a-z, 0-9, _) " +
				"Please try entering another username!",
		})
		return false
	}
	// Duplicate username check
	cur, _ := usercollection.Find(context.TODO(), bson.M{})
	for cur.Next(context.TODO()) {
		// create a value into which the single document can be decoded
		var elem models.UserModel
		cur.Decode(&elem)
		if elem.Username == username {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Sorry, this username has already been taken. Please select another username and try again!",
			})
			return false
		}
	}

	return true
}

// ValidateEmail checks whether the provided email is a valid email address
// This includes banning commonly banned characters, checking for the @ character, and limiting email address length (to 75 characters)
func ValidateEmail(c *gin.Context, email string) bool {
	// Invalid email check
	emailRE := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_\x60{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	emailCharLimit := 75
	if len(email) > emailCharLimit {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": fmt.Sprintf("This email is too long, please select an email address that is less than %d characters long and try again.", emailCharLimit),
		})
		return false
	} else if !emailRE.MatchString(email) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "This email is invalid. Please make sure to check whether you typed your email address correctly or use a different email address and try again.",
		})
		return false
	}

	return true
}

// ValidateURL checks whether the given URL is valid.
// It uses govalidator's IsURL method.
func ValidateURL(c *gin.Context, URL string) bool {
	// Invalid URL check
	if !govalidator.IsURL(URL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "This URL is invalid. Please check to see whether your URL entry is correct and try again.",
		})
		return false
	}

	return true
}
