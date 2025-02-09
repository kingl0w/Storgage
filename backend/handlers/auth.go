package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"storgage/config"
	"storgage/database"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Invite   string `json:"invite,omitempty"`
}

// SignUp example using jsonError from helpers.go
func Signup(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		jsonError(w, "Invalid request format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if creds.Username == "" || creds.Password == "" || creds.Invite == "" {
		jsonError(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Check if username already exists
	var exists bool
	err = database.DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)",
		creds.Username).Scan(&exists)
	if err != nil {
		jsonError(w, "Error checking username", http.StatusInternalServerError)
		return
	}
	if exists {
		jsonError(w, "Username already taken", http.StatusConflict)
		return
	}

	// Check if invite code is valid and not used
	var inviteUsed bool
	err = database.DB.QueryRow(context.Background(),
		"SELECT used FROM invite_codes WHERE code = $1",
		creds.Invite).Scan(&inviteUsed)
	if err != nil {
		// If the code wasn't found in DB, that's also an error:
		jsonError(w, "Invalid invite code", http.StatusForbidden)
		return
	}
	if inviteUsed {
		jsonError(w, "Invite code already used", http.StatusForbidden)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		jsonError(w, "Error processing password", http.StatusInternalServerError)
		return
	}

	// Begin transaction
	tx, err := database.DB.Begin(context.Background())
	if err != nil {
		jsonError(w, "Error starting transaction", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback(context.Background())

	// Store new user
	var userId int
	err = tx.QueryRow(context.Background(),
		"INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id",
		creds.Username, string(hashedPassword)).Scan(&userId)
	if err != nil {
		jsonError(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Mark invite as used
	_, err = tx.Exec(context.Background(),
		"UPDATE invite_codes SET used = TRUE, used_by = $1, used_at = $2 WHERE code = $3",
		userId, time.Now(), creds.Invite)
	if err != nil {
		jsonError(w, "Error updating invite code", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err = tx.Commit(context.Background()); err != nil {
		jsonError(w, "Error completing registration", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// Login is unchanged except for using jsonError
func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		jsonError(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var hashedPassword string
	err = database.DB.QueryRow(context.Background(),
		"SELECT password FROM users WHERE username = $1",
		creds.Username).Scan(&hashedPassword)
	if err != nil {
		jsonError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		jsonError(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": creds.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.LoadConfig().JWTSecret))
	if err != nil {
		jsonError(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
