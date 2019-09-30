/*
Package routes contains CRUD-to-REST mapping, making the code more accessible through an API.
*/
package routes

import (
	"fir/controllers"

	"github.com/gin-gonic/gin"
)

// InitUserRoutes initializes (routes) the operations concerning users. This includes all CRUD operations.
func InitUserRoutes(routerGroup *gin.RouterGroup) {
	userBaseURL := "/users"
	urgroup := routerGroup.Group(userBaseURL)
	{
		urgroup.POST("", controllers.PostUser)
		urgroup.GET("", controllers.GetAllUsers)
		urgroup.GET("/:userID", controllers.GetUser)
		urgroup.PUT("/:userID", controllers.UpdateUser)
		urgroup.DELETE("/:userID", controllers.DeleteUser)
	}
}

// InitIncidentRoutes initializes (routes) the operations concerning incidents. This includes all CRUD operations.
func InitIncidentRoutes(routerGroup *gin.RouterGroup) {
	incidentBaseURL := "/incidents"
	irgroup := routerGroup.Group(incidentBaseURL)
	{
		irgroup.POST("", controllers.PostIncident)
		irgroup.GET("", controllers.GetAllIncidents)
		irgroup.GET("/:incidentID", controllers.GetIncident)
		irgroup.PUT("/:incidentID", controllers.UpdateIncident)
		irgroup.DELETE("/:incidentID", controllers.DeleteIncident)
	}
}

// InitArtifactRoutes initializes (routes) the operations concerning Artifacts. This includes all CRUD operations.
func InitArtifactRoutes(routerGroup *gin.RouterGroup) {
	artifactBaseURL := "/artifacts"
	argroup := routerGroup.Group(artifactBaseURL)
	{
		argroup.POST("", controllers.PostArtifact)
		argroup.GET("", controllers.GetAllArtifacts)
		argroup.GET("/:artifactID", controllers.GetArtifact)
		argroup.PUT("/:artifactID", controllers.UpdateArtifact)
		argroup.DELETE("/:artifactID", controllers.DeleteArtifact)
	}
}

// InitFileRoutes initializes (routes) the operations concerning Files. This includes all CRUD operations.
func InitFileRoutes(routerGroup *gin.RouterGroup) {
	fileBaseURL := "/files"
	frgroup := routerGroup.Group(fileBaseURL)
	{
		frgroup.POST("", controllers.PostFile)
		frgroup.GET("", controllers.GetAllFiles)
		frgroup.GET("/:fileID", controllers.GetFile)
		frgroup.PUT("/:fileID", controllers.UpdateFile)
		frgroup.DELETE("/:fileID", controllers.DeleteFile)
	}
}
