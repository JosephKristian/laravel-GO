package http

import (
	"encoding/json"
	"net/http"

	"github.com/JosephKristian/project-migration/internal/models"
	"github.com/go-playground/validator/v10"
)

type RegisterController struct {
	registerService RegisterService
	validate        *validator.Validate
}

func (r *RegisterController) Register(w http.ResponseWriter, req *http.Request) {
	var user models.User

	// Decode JSON request body to User struct
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate input fields
	if err := r.validate.Struct(&user); err != nil {
		http.Error(w, "Validation failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Process registration
	registeredUser, err := r.registerService.Register(user)
	if err != nil {
		http.Error(w, "Registration failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registeredUser)
}
