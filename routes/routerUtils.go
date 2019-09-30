package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thinkerou/favicon"
)

// SetupRouter initializes a gin engine with all of the routes and middleware
func SetupRouter() *gin.Engine {
	// Initialize gin router
	router := gin.Default()

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	router.Use(favicon.New("img/favicon.ico")) // set favicon middleware

	// Initialize routing groups for all methods
	routerGroup := router.Group("/api")
	{
		InitAuthRoutes(routerGroup)
		InitUserRoutes(routerGroup)
		InitIncidentRoutes(routerGroup)
		InitArtifactRoutes(routerGroup)
		InitFileRoutes(routerGroup)
	}

	return router
}
