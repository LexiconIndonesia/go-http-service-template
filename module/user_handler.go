package module

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/adryanev/go-http-service-template/common/models"
	"github.com/adryanev/go-http-service-template/common/utils"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// userHandler handles user-related requests
func (m *Module) userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		m.createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// createUser creates a new user
func (m *Module) createUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		log.Error().Err(err).Msg("Failed to decode request body")
		utils.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Generate UUID for new user
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	// Validate user
	validator := models.NewUserValidator()
	if err := validator.Validate(user); err != nil {
		log.Error().Err(err).Msg("Invalid user data")
		utils.WriteError(w, http.StatusBadRequest, "Invalid user data: "+err.Error())
		return
	}

	// TODO: Save user to database
	// In a real implementation, you would save the user to the database here
	// using the dependency-injected DB

	// Return success
	utils.WriteJSON(w, http.StatusCreated, user)
}
