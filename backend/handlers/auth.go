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

// Signup handles user registration
func Signup(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Check if invite code is valid and not used
	var inviteUsed bool
	err = database.DB.QueryRow(context.Background(), "SELECT used FROM invite_codes WHERE code = $1", creds.Invite).Scan(&inviteUsed)
	if err != nil || inviteUsed {
		http.Error(w, "Invalid or already used invite code", http.StatusForbidden)
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(creds.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	// Store new user
	_, err = database.DB.Exec(context.Background(), "INSERT INTO users (username, password) VALUES ($1, $2)", creds.Username, string(hashedPassword))
	if err != nil {
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	// Mark invite as used
	_, _ = database.DB.Exec(context.Background(), "UPDATE invite_codes SET used = TRUE WHERE code = $1", creds.Invite)

	w.WriteHeader(http.StatusCreated)
}

// Login handles user authentication
func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	var hashedPassword string
	err = database.DB.QueryRow(context.Background(), "SELECT password FROM users WHERE username = $1", creds.Username).Scan(&hashedPassword)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": creds.Username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(config.LoadConfig().JWTSecret))
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Send token
	json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
}
