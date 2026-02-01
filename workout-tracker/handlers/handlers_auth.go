package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"workout-tracker/auth"
	"workout-tracker/database"
	"workout-tracker/middleware"
	"workout-tracker/models"
)

// RegisterRequest represents the data needed to register a new user
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the data needed to login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse is what we send back after successful login
type LoginResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// Register creates a new user account
func Register(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Username, email, and password are required", http.StatusBadRequest)
		return
	}

	// Hash the password before storing it
	passwordHash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Insert the new user into the database
	var userID int
	err = database.DB.QueryRow(
		`INSERT INTO users (username, email, password_hash) 
		 VALUES ($1, $2, $3) RETURNING id`,
		req.Username, req.Email, passwordHash,
	).Scan(&userID)

	if err != nil {
		// Check if it's a duplicate username/email error
		if err.Error() == "pq: duplicate key value violates unique constraint \"users_username_key\"" ||
			err.Error() == "pq: duplicate key value violates unique constraint \"users_email_key\"" {
			http.Error(w, "Username or email already exists", http.StatusConflict)
			return
		}
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Return the created user (without password hash)
	user := models.User{
		ID:       userID,
		Username: req.Username,
		Email:    req.Email,
	}

	// Generate JWT token for the new user
	token, err := auth.GenerateToken(userID, req.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Send response with token and user
	response := LoginResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// Login authenticates a user and returns a JWT token
func Login(w http.ResponseWriter, r *http.Request) {
	// Parse the request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Username == "" || req.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	// Find the user in the database
	var user models.User
	var passwordHash string
	err := database.DB.QueryRow(
		`SELECT id, username, email, password_hash, created_at 
		 FROM users WHERE username = $1`,
		req.Username,
	).Scan(&user.ID, &user.Username, &user.Email, &passwordHash, &user.CreatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Error finding user", http.StatusInternalServerError)
		return
	}

	// Check if the password is correct
	if !auth.CheckPassword(req.Password, passwordHash) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate a JWT token
	token, err := auth.GenerateToken(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Return the token and user info
	response := LoginResponse{
		Token: token,
		User:  user,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCurrentUser returns information about the currently logged-in user
func GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user info from the context (set by AuthMiddleware)
	claims := middleware.GetUserFromContext(r)
	if claims == nil {
		http.Error(w, "User not found in context", http.StatusUnauthorized)
		return
	}

	// Fetch full user details from database
	var user models.User
	err := database.DB.QueryRow(
		`SELECT id, username, email, created_at FROM users WHERE id = $1`,
		claims.UserID,
	).Scan(&user.ID, &user.Username, &user.Email, &user.CreatedAt)

	if err != nil {
		http.Error(w, "Error fetching user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
