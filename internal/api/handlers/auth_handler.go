package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"src/internal/data"
	"src/internal/helpers"
)

type MyHandler struct {
	DB *sql.DB
}

func (h MyHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var newUser data.UserRegister
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Hash the user's password
	if newUser.Password != newUser.ConfirmPassword {
		http.Error(w, "Password did not matched", http.StatusBadRequest)
	}
	hashedPassword, err := helpers.HashPassword(newUser.Password)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	// Insert user into the database with the hashed password
	_, err = h.DB.Exec("INSERT INTO users (username, password_hash, role) VALUES ($1, $2, $3)",
		newUser.Username, hashedPassword, "user")
	if err != nil {
		http.Error(w, "Failed to insert user", http.StatusInternalServerError)
		return
	}

	// Respond with a success message or status code
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User registration successful")
}

func (h MyHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	// Parse JSON request body
	var loginRequest data.UserLogin
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(w, "Invalid request data", http.StatusBadRequest)
		return
	}

	// Query the database to retrieve the user's hashed password
	var user data.User
	err = h.DB.QueryRow("SELECT id, username, password_hash, role FROM users WHERE username = $1", loginRequest.Username).Scan(
		&user.ID,
		&user.Username,
		&user.PasswordHash,
		&user.Role,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "User not found", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Verify the provided password against the stored hash
	if err := helpers.VerifyPassword(user.PasswordHash, loginRequest.Password); err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Authentication successful
	// Generate a JWT token
	tokenString, refreshToken, err := helpers.GenerateTokens(user.ID, user.Username)
	if err != nil {
		http.Error(w, "Failed to generate JWT token", http.StatusInternalServerError)
		return
	}

	// Include the token in the response
	response := map[string]string{"token": tokenString, "refresh_token": refreshToken}

	// Respond with the token and a success status code
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
