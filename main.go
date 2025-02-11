package main

import (
	"log"
	"os"

	"github.com/dro14/nuqta-service/handler"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

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

	file, err = os.Create("my.log")
	if err != nil {
		log.Fatal("can't open my.log: ", err)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	h := handler.New()
	h.UpdateRecs()
	err = h.Run(port)
	if err != nil {
		log.Fatal("Error running handler: ", err)
	}
}
