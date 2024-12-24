package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"firebase.google.com/go/v4/auth"
	"github.com/dro14/nuqta-service/database"
	"github.com/dro14/nuqta-service/repository"
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

func CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request map[string]string
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username, ok := request["username"]
	if !ok {
		http.Error(w, "Missing username", http.StatusBadRequest)
		return
	}
	bio, ok := request["bio"]
	if !ok {
		http.Error(w, "Missing bio", http.StatusBadRequest)
		return
	}

	err = repository.CreateUser(r.Context(), username, bio)
	if err != nil {
		log.Print("Failed to create user: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"message": "User created"})
	w.Header().Set("Content-Type", "application/json")
}

func FollowUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request map[string]string
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	follower, ok := request["follower"]
	if !ok {
		http.Error(w, "Missing follower", http.StatusBadRequest)
		return
	}
	followee, ok := request["followee"]
	if !ok {
		http.Error(w, "Missing followee", http.StatusBadRequest)
		return
	}

	err = repository.FollowUser(r.Context(), follower, followee)
	if err != nil {
		log.Print("Failed to follow user: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"message": "User followed"})
	w.Header().Set("Content-Type", "application/json")
}

func CreatePost(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request map[string]string
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username, ok := request["username"]
	if !ok {
		http.Error(w, "Missing username", http.StatusBadRequest)
		return
	}
	content, ok := request["content"]
	if !ok {
		http.Error(w, "Missing content", http.StatusBadRequest)
		return
	}

	err = repository.CreatePost(r.Context(), username, content)
	if err != nil {
		log.Print("Failed to create post: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"message": "Post created"})
	w.Header().Set("Content-Type", "application/json")
}

func GetFeed(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var request map[string]string
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	username, ok := request["username"]
	if !ok {
		http.Error(w, "Missing username", http.StatusBadRequest)
		return
	}

	feed, err := repository.GetFeed(r.Context(), username)
	if err != nil {
		log.Print("Failed to get feed: ", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(feed)
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

	err = database.Init(os.Getenv("NEO4J_URI"), "neo4j", os.Getenv("NEO4J_PASSWORD"))
	if err != nil {
		log.Fatal("Failed to connect to Neo4j: ", err)
	}

	http.HandleFunc("/", Root)
	http.HandleFunc("/auth", firebaseAuth.AuthMiddleware(Auth))
	http.HandleFunc("/create_user", CreateUser)
	http.HandleFunc("/follow_user", FollowUser)
	http.HandleFunc("/create_post", CreatePost)
	http.HandleFunc("/get_feed", GetFeed)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
