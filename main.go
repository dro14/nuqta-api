package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"firebase.google.com/go/v4/auth"
	"github.com/joho/godotenv"
)

func RootHandler(w http.ResponseWriter, _ *http.Request) {
	response := map[string]string{
		"message": "Hello, World!",
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func AuthHandler(w http.ResponseWriter, r *http.Request) {
	token, ok := r.Context().Value("token").(*auth.Token)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	response := map[string]string{
		"uid":   token.UID,
		"email": token.Claims["email"].(string),
	}

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("can't load .env file: ", err)
	}

	firebaseAuth, err := NewFirebaseAuth("service_account_key.json")
	if err != nil {
		log.Fatal("Error initializing Firebase: ", err)
	}

	http.HandleFunc("/", RootHandler)
	http.HandleFunc("/auth", firebaseAuth.AuthMiddleware(AuthHandler))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
