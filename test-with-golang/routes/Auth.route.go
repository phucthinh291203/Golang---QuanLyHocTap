package routes

import (
	controllers "test-with-golang/Controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(router *gin.Engine) {
	authAdminGroup := router.Group("/auth")
	{
		authAdminGroup.POST("/login", func(ctx *gin.Context) {
			controllers.Login(ctx)
		})
	}
}



