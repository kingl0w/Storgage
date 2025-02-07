package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"storgage/database"
	"time"
)

type InviteResponse struct {
	Code string `json:"code"`
}

// Generate a random invite code
func generateInviteCode() string {
	rand.Seed(time.Now().UnixNano())
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 8) // 8-character invite code
	for i := range code {
		code[i] = charset[rand.Intn(len(charset))]
	}
	return string(code)
}

// Create a new invite code and store it in the database
func GenerateInvite(w http.ResponseWriter, r *http.Request) {
	inviteCode := generateInviteCode()

	// Store the invite code in the database
	_, err := database.DB.Exec(context.Background(), "INSERT INTO invite_codes (code, used) VALUES ($1, $2)", inviteCode, false)
	if err != nil {
		http.Error(w, "Failed to create invite code", http.StatusInternalServerError)
		return
	}

	// Return the invite code
	response := InviteResponse{Code: inviteCode}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
