package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	repo *UserRepository
}

// NewUserHandler creates a new user handler
func NewUserHandler(repo *UserRepository) *UserHandler {
	return &UserHandler{repo: repo}
}

// CreateUser handles POST /users
// CreateUser handles POST /users
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User to create"
// @Success 201 {object} User
// @Failure 400 {string} string "Invalid request body or missing fields"
// @Failure 500 {string} string "Error creating user"
// @Router /users [post]
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Username == "" || req.Name == "" || req.Email == "" || req.Age <= 0 {
		http.Error(w, "Username, name, email, and age are required", http.StatusBadRequest)
		return
	}

	user, err := h.repo.CreateUser(req)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// GetUser handles GET /users/{username}
// GetUser handles GET /users/{username}
// @Summary Get a user by username
// @Description Retrieve a user's details by their username
// @Tags users
// @Produce json
// @Param username path string true "Username of the user to retrieve"
// @Success 200 {object} User
// @Failure 400 {string} string "Invalid username"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Error getting user"
// @Router /users/{username} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract username from path
	username, err := extractUsernameFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	user, err := h.repo.GetUserByUsername(username)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting user: %v", err)
		http.Error(w, "Error getting user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// GetAllUsers handles GET /users
// GetAllUsers handles GET /users
// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags users
// @Produce json
// @Success 200 {array} User
// @Failure 500 {string} string "Error getting users"
// @Router /users [get]
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.repo.GetAllUsers()
	if err != nil {
		log.Printf("Error getting users: %v", err)
		http.Error(w, "Error getting users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// UpdateUser handles PUT /users/{username}
// UpdateUser handles PUT /users/{username}
// @Summary Update an existing user
// @Description Update an existing user's details by their username
// @Tags users
// @Accept json
// @Produce json
// @Param username path string true "Username of the user to update"
// @Param user body UpdateUserRequest true "User update details"
// @Success 200 {object} User
// @Failure 400 {string} string "Invalid username or request body"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Error updating user"
// @Router /users/{username} [put]
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract username from path
	username, err := extractUsernameFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.Name == "" || req.Email == "" || req.Age <= 0 {
		http.Error(w, "Name, email, and age are required", http.StatusBadRequest)
		return
	}

	user, err := h.repo.UpdateUser(username, req)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Error updating user: %v", err)
		http.Error(w, "Error updating user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// DeleteUser handles DELETE /users/{username}
// DeleteUser handles DELETE /users/{username}
// @Summary Delete a user
// @Description Delete a user by their username
// @Tags users
// @Param username path string true "Username of the user to delete"
// @Success 204 "No Content"
// @Failure 400 {string} string "Invalid username"
// @Failure 404 {string} string "User not found"
// @Failure 500 {string} string "Error deleting user"
// @Router /users/{username} [delete]
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract username from path
	username, err := extractUsernameFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	err = h.repo.DeleteUser(username)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Error deleting user: %v", err)
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// extractUsernameFromPath extracts the username from the URL path
// Expects path format: /users/{username}
func extractUsernameFromPath(path string) (string, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) < 2 {
		return "", http.ErrNotSupported
	}

	username := parts[1]
	if username == "" {
		return "", http.ErrNotSupported
	}

	return username, nil
}
