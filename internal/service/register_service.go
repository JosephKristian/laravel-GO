package service

import (
	"errors"

	"github.com/username/project-migration/internal/models"
)

type RegisterService struct{}

func (r *RegisterService) Register(user models.User) (*models.User, error) {
	// Implement the business logic for registration, such as saving to the database
	// For now, just return the user as a mock successful registration
	if user.Name == "" || user.Email == "" || user.Phone == "" || user.Password == "" {
		return nil, errors.New("all fields are required")
	}

	// Normally you would interact with a database here

	return &user, nil
}
