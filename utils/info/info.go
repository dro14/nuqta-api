package info

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const adminUserID = 1331278972

var bot *tgbotapi.BotAPI

func SetUp() {
	token, ok := os.LookupEnv("INFO_BOT_TOKEN")
	if !ok {
		log.Fatal("info bot token is not specified")
	}

	var err error
	bot, err = tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("can't initialize info bot: ", err)
	}

	go ReceiveUpdates()
}

func ReceiveUpdates() {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates := bot.GetUpdatesChan(updateConfig)
	for update := range updates {
		message := update.Message
		if message != nil && message.From.ID == adminUserID {
			switch message.Command() {
			case "logs":
				SendDocument("my.log")
				SendDocument("gin.log")
			case "crashed_logs":
				SendDocument("my-crashed.log")
				SendDocument("gin-crashed.log")
			}
		}
	}
}

func MonitorShutdown(sigChan chan os.Signal) {
	sig := <-sigChan
	log.Printf("Received %v, initiating shutdown...", sig)
	log.SetOutput(os.Stdout)
	gin.DefaultWriter = os.Stdout
	SendDocument("my.log")
	SendDocument("gin.log")
}

func SendMessage(text string) {
	config := tgbotapi.NewMessage(adminUserID, text)
	_, err := bot.Request(config)
	if err != nil {
		log.Print("can't send info message: ", err)
	}
}

func SendDocument(path string) {
	config := tgbotapi.NewDocument(adminUserID, tgbotapi.FilePath(path))
	_, err := bot.Request(config)
	if err != nil {
		log.Print("can't send info document: ", err)
	}
}
