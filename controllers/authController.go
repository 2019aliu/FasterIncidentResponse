package controllers

import (
	"context"
	"fir/models"
	mongoFunctions "fir/mongo"
	redisFunctions "fir/redis"
	"net/http"

	"github.com/gin-gonic/gin"
	uuid "github.com/nu7hatch/gouuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

// Signin is the function used to sign users into the FaIR system
func Signin(c *gin.Context) {
	var incomingUser models.UserModel
	UnmarshalRequest(c, &incomingUser)
	var usercollection = mongoFunctions.GetCollection(mongoFunctions.GetClient(), "users")

	var matchedUser models.UserModel
	err := usercollection.FindOne(context.TODO(), bson.M{"username": incomingUser.Username}).Decode(&matchedUser)
	if err != nil {
		c.String(http.StatusNotFound, "This username was not found")
		// log.Fatal(err)
		return
	}

	// CompareHashAndPassword returns error if the passwords don't match, also put the hashed one (stored) first)
	if err = bcrypt.CompareHashAndPassword([]byte(matchedUser.Password), []byte(incomingUser.Password)); err != nil {
		c.String(http.StatusUnauthorized, "The password provided is an incorrect password for this username")
		return
	}

	// Create a new random session token
	newUUID, err := uuid.NewV4()
	sessionToken := newUUID.String()
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	var cache = redisFunctions.InitCache()
	_, err = cache.Do("SETEX", sessionToken, "120", incomingUser.Username)
	if err != nil {
		// If there is an error in setting the cache, return an internal server error
		c.String(http.StatusInternalServerError, "There was an error setting up the token in the redis cache")
		return
	}

	// expire := time.Now().Add(120*time.Second).Unix()
	cookie := http.Cookie{
		Name:   "sessiontoken",
		Value:  sessionToken,
		MaxAge: 120,
		Path:   "/",
		// Domain:   "localhost:8080",
		Secure:   true,
		HttpOnly: true,
	}

	var w http.ResponseWriter = c.Writer
	// var req *http.Request = c.Request

	http.SetCookie(w, &cookie)
	// c.SetCookie("session_token", sessionToken, int(time.Now().Add(120*time.Second).Unix()), "/api", "localhost", true, true)

	// cookie, err := c.Request.Cookie("session_token")
	// if err != nil {
	// 	if err == http.ErrNoCookie {
	// 		// If the cookie is not set, return an unauthorized status
	// 		c.String(http.StatusUnauthorized, "The cookie was not set")
	// 		return
	// 	}
	// 	// For any other type of error, return a bad request status
	// 	c.String(http.StatusBadRequest, "Something happened in between :( please try again")
	// 	return
	// }
	// fmt.Println(cookie)
	// sessionToken = cookie.Value

	c.JSON(http.StatusOK, matchedUser)
}

// func GetAllTokens
