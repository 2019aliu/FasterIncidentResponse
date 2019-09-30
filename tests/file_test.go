package tests

import (
	"encoding/json"
	"fir/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"fir/routes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPostfile(t *testing.T) {
	// Initialize gin router
	router := gin.Default()
	router.Use(gin.Recovery())
	// routes.InitFileRoutes(router)

	// Use net/http/httptest's testing package
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/files", strings.NewReader(`{
		"filepath": "thisisanotherfile.txt"
	}`))
	router.ServeHTTP(w, req)

	assert.Equal(t, 201, w.Code)
	assert.Equal(t, "hi", w.Body.String())
}

func TestGetFile(t *testing.T) {
	// Get the router
	router := routes.SetupRouter()

	// Use net/http/httptest's testing package
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/files/5d3a2ed1ab6a0fe1eb05f95f", nil)
	router.ServeHTTP(w, req)

	// Basic assertions with testify/assert's extension in Golang
	assert.Equal(t, 200, w.Code)
	var testFile models.FileModel
	json.Unmarshal([]byte(w.Body.String()), &testFile)
	assert.Equal(t, "/home/aliu/testing.txt", testFile.FilePath)
}
