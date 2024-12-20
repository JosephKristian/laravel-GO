package http

import (
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine, registerController *RegisterController) {
	// Menggunakan /db/v1 sebagai prefix untuk semua rute API
	api := router.Group("/db/v1")
	{
		// Grup rute untuk otentikasi
		auth := api.Group("/auth")
		{
			// Route untuk registrasi pengguna
			auth.POST("/register", registerController.Register)

			// Route untuk aktivasi akun pengguna
			auth.POST("/account-activation", registerController.AccountActivation)
		}
	}
}
