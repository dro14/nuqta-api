package main

import (
	"log"
	"os"

	"github.com/dro14/nuqta-service/handler"
	"github.com/dro14/nuqta-service/utils/info"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("can't load .env file: ", err)
	}
	info.SetUp()

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	info.SendMessage("Nuqta service restarted")
	h := handler.New()
	err = h.Run(port)
	if err != nil {
		log.Fatal("Error running handler: ", err)
	}
}
