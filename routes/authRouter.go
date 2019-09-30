package routes

import (
	"fir/controllers"

	"github.com/gin-gonic/gin"
)

// InitAuthRoutes initializes all routes related to authenticating users
func InitAuthRoutes(router *gin.RouterGroup) {
	authBaseURL := "/"
	authRouter := router.Group(authBaseURL)
	{
		authRouter.POST("/signin", controllers.Signin)
		authRouter.POST("/signup", controllers.PostUser)
	}
}

// func InitTokenRoutes(router *gin.RouterGroup) {
// 	tokenBaseURL := "/token"
// 	trgroup := router.Group(tokenBaseURL)
// 	{
// 		trgroup.POST("", controllers.PostToken)
// 		trgroup.GET("", controllers.GetAllTokens)
// 		trgroup.GET("/:userID", controllers.GetToken)
// 		trgroup.PUT("/:userID", controllers.UpdateToken)
// 		trgroup.DELETE("/:userID", controllers.DeleteToken)
// 	}
// }
