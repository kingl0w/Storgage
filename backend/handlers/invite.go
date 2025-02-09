package handlers

import (
	"context"
	"encoding/json"
	"math/rand"
	"net/http"
	"os"
	"storgage/database"
	"time"
)

type InviteResponse struct {
	Code string `json:"code"`
}

type AdminCredentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// generate invite code
func generateInviteCode() string {
	// Use crypto/rand in production for better randomness
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 8)
	for i := range code {
		code[i] = charset[rnd.Intn(len(charset))]
	}
	return string(code)
}

// create invite code and store it in database
func GenerateInvite(w http.ResponseWriter, r *http.Request) {
	// Parse admin credentials
	var creds AdminCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//verify admin credentials
	if creds.Username != os.Getenv("ADMIN_USERNAME") ||
		creds.Password != os.Getenv("ADMIN_PASSWORD") {
		http.Error(w, "Invalid admin credentials", http.StatusUnauthorized)
		return
	}

	inviteCode := generateInviteCode()

	//store invite in the database
	_, err := database.DB.Exec(
		context.Background(),
		`INSERT INTO invite_codes (code, used, created_at) 
         VALUES ($1, $2, $3)`,
		inviteCode,
		false,
		time.Now(),
	)

	if err != nil {
		http.Error(w, "Failed to create invite code", http.StatusInternalServerError)
		return
	}

	//return invite code
	response := InviteResponse{Code: inviteCode}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// verify invite code when user signs up
func VerifyInviteCode(code string) (bool, error) {
	var used bool
	err := database.DB.QueryRow(
		context.Background(),
		`SELECT used FROM invite_codes WHERE code = $1`,
		code,
	).Scan(&used)

	if err != nil {
		return false, err
	}

	if used {
		return false, nil
	}

	//mark code as used
	_, err = database.DB.Exec(
		context.Background(),
		`UPDATE invite_codes SET used = true, used_at = $1 WHERE code = $2`,
		time.Now(),
		code,
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

// VerifyInviteHandler handles checking the validity of an invite code
func VerifyInviteHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	valid, err := VerifyInviteCode(request.Code)
	if err != nil {
		http.Error(w, "Error verifying invite code", http.StatusInternalServerError)
		return
	}

	if !valid {
		http.Error(w, "Invite code is invalid or already used", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Invite code is valid!"})
}
