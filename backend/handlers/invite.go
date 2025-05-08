package handlers

import (
	"context"
	"encoding/json"
	"log"
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

func generateInviteCode() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	code := make([]byte, 8)
	for i := range code {
		code[i] = charset[rnd.Intn(len(charset))]
	}
	return string(code)
}

func GenerateInvite(w http.ResponseWriter, r *http.Request) {
	var creds AdminCredentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	//verify admin credentials
	if creds.Username != os.Getenv("ADMIN_USERNAME") ||
		creds.Password != os.Getenv("ADMIN_PASSWORD") {
		jsonError(w, "Invalid admin credentials", http.StatusUnauthorized)
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
		jsonError(w, "Failed to create invite code", http.StatusInternalServerError)
		return
	}

	response := InviteResponse{Code: inviteCode}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

//define the VerifyInviteCode function if needed
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
	return true, nil
}

//verifyInviteHandler calls verifyInviteCode
func VerifyInviteHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received verify-invite request from: %s", r.RemoteAddr)
	var request struct {
		Code string `json:"code"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		jsonError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	valid, err := VerifyInviteCode(request.Code)
	if err != nil {
		jsonError(w, "Error verifying invite code", http.StatusInternalServerError)
		return
	}

	if !valid {
		jsonError(w, "Invite code is invalid or already used", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Invite code is valid!"})
}
