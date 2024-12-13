package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
	"github.com/username/project-migration/internal/delivery/http"
	"github.com/username/project-migration/internal/service"
)

func main() {
	// Create a Gin router
	router := gin.Default()

	// Initialize services and controllers
	registerService := &service.RegisterService{}
	registerController := &http.RegisterController{
		registerService: registerService,
		validate:        validator.New(),
	}

	// Register routes
	http.RegisterRoutes(router)

	// Start the server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
