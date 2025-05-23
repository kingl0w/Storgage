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

	r := mux.NewRouter()

	//logging middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s", r.Method, r.URL.Path)
			next.ServeHTTP(w, r)
		})
	})

	//CORS 
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, r)
		})
	})

	//routes 
	r.HandleFunc("/api/files", storageHandler.ListFiles).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/files/{filename}", storageHandler.DeleteFile).
		Methods("DELETE", "OPTIONS")
	r.HandleFunc("/api/upload", storageHandler.UploadFile).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/signup", handlers.Signup).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/login", handlers.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/admin/invite", handlers.GenerateInvite).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/verify-invite", handlers.VerifyInviteHandler).Methods("POST", "OPTIONS")

	port := config.Port
	log.Printf("Server starting on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
