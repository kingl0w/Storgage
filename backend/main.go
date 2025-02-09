package main

import (
	"log"
	"net/http"
	"storgage/config"
	"storgage/database"
	"storgage/handlers"

	"github.com/gorilla/mux"
)

func main() {
	config := config.LoadConfig()
	database.ConnectDB(config)

	storageHandler, err := handlers.NewStorageHandler(config)
	if err != nil {
		log.Fatal(err)
	}

	//initialize router
	r := mux.NewRouter()

	//CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Allow any origin during development
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			// Handle preflight requests
			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	//routes
	r.HandleFunc("/api/files", storageHandler.ListFiles).Methods("GET")
	r.HandleFunc("/api/upload", storageHandler.UploadFile).Methods("POST")
	r.HandleFunc("/api/signup", handlers.Signup).Methods("POST")
	r.HandleFunc("/api/login", handlers.Login).Methods("POST")
	r.HandleFunc("/api/admin/invite", handlers.GenerateInvite).Methods("POST")
	r.HandleFunc("/api/verify-invite", handlers.VerifyInviteHandler).Methods("POST")

	//start server
	port := config.Port
	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))

}
