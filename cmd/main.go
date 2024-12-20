package main

import (
	"log"

	"github.com/JosephKristian/project-migration/internal/delivery/http"
	"github.com/JosephKristian/project-migration/internal/repositories"
	"github.com/JosephKristian/project-migration/internal/service"
	"github.com/JosephKristian/project-migration/internal/usecase/database" // Import package database
	"github.com/gin-gonic/gin"
)

func main() {
	// Inisialisasi database menggunakan fungsi dari internal/database/db.go
	db, err := database.InitDB() // Panggil fungsi InitDB dari package database
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Inisialisasi repository dengan database
	otpRepository := &repositories.OtpRepo{DB: db}
	otpService := service.NewOtpService(otpRepository)

	// Inisialisasi user repository dan service
	userRepo := repositories.UserRepo{DB: db}
	registerService := service.NewRegisterService(db, userRepo, otpService) // Tambahkan otpService
	accountActivationService := service.NewAccountActivationService()

	// Inisialisasi controller dengan dependency injection
	registerController := http.NewRegisterController(registerService, accountActivationService)

	// Inisialisasi router
	router := gin.Default()

	// Daftarkan routes
	http.RegisterRoutes(router, registerController)

	// Start server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
