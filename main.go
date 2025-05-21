package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/dro14/nuqta-service/handler"
	"github.com/dro14/nuqta-service/utils/info"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	file, err := os.Create("my.log")
	if err != nil {
		log.Fatal("can't open my.log: ", err)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	file, err = os.Create("gin.log")
	if err != nil {
		log.Fatal("can't open gin.log: ", err)
	}
	gin.DefaultWriter = file
	gin.SetMode(gin.ReleaseMode)

	err = godotenv.Load()
	if err != nil {
		log.Fatal("can't load .env file: ", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go info.MonitorShutdown(sigChan)

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
