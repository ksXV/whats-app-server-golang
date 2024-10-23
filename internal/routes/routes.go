package routes

import (
	"whatsapp-server/internal/auth"
	"whatsapp-server/internal/database"
	"whatsapp-server/internal/oauth"

	"github.com/gin-gonic/gin"
)

func GroupRoutes(router *gin.Engine, dbConn database.DBConnection) {
	api := router.Group("/api")
	{
		groupAuthRoutes(api, dbConn)
	}
}

func groupAuthRoutes(group *gin.RouterGroup, dbConn database.DBConnection) {
	authGroup := group.Group("/auth")
	{
		authGroup.GET("/login/google", oauth.HandleGoogleAuthLink)
		authGroup.GET("/callback/google", oauth.HandleGoogleCallBack(dbConn))

		authGroup.GET("/login", auth.HandleLogin(dbConn))
		authGroup.POST("/register", auth.HandleRegister)
		authGroup.POST("/forgot-password", nil)
	}
}
