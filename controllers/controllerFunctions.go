package controllers

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

// UnmarshalRequest unmarshals the request body, which is assumed to be a proper json.
// This method is a helper method for any method that requires parsing the request body (currently post and put operations)
func UnmarshalRequest(c *gin.Context, unmarshalObj interface{}) []byte {
	buf := make([]byte, 2048) // 2048 is an arbitrary number whatever
	num, _ := c.Request.Body.Read(buf)
	reqBody := buf[0:num] // don't cast to string otherwise unmarshal won't work
	json.Unmarshal(reqBody, unmarshalObj)
	return reqBody
}
