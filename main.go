package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"firebase.google.com/go/v4/auth"
	"github.com/dro14/nuqta-service/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Root(w http.ResponseWriter, _ *http.Request) {
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

func Auth(w http.ResponseWriter, r *http.Request) {
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

	file, err := os.Create("gin.log")
	if err != nil {
		log.Fatal("can't open gin.log: ", err)
	}
	gin.DefaultWriter = file

	file, err = os.Create("nuqta-service.log")
	if err != nil {
		log.Fatal("can't open yordamchi.log: ", err)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// _, err = NewFirebaseAuth("service_account_key.json")
	// if err != nil {
	//     log.Fatal("Error initializing Firebase: ", err)
	// }

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	h := handler.New()
	err = h.Run(port)
	if err != nil {
		log.Fatal("Error running handler: ", err)
	}
}
