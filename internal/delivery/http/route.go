package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api := router.Group("/v1")
	{
		auth := api.Group("/auth")
		{
			auth.POST("/register", Register)
		}
	}
}
