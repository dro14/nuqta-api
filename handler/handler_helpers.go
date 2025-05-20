package handler

import (
	"log"
	"net/http"
	"runtime/debug"

	"github.com/dro14/nuqta-service/models"
	"github.com/dro14/nuqta-service/utils/info"
	"github.com/gin-gonic/gin"
)

func failure(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func sendSSEEvent(c *gin.Context, name string, data any) {
	c.SSEvent(name, data)
	c.Writer.Flush()
}

func broadcast(uid string, name string, data any) {
	broadcastersMutex.RLock()
	channels := broadcasters[uid]
	for _, channel := range channels {
		channel <- &models.Event{Name: name, Data: data}
	}
	broadcastersMutex.RUnlock()
}

func notifyOnPanic(c *gin.Context, err any) {
	log.Printf("%s\n%s", err, debug.Stack())
	info.SendDocument("gin.log")
	info.SendDocument("yordamchi.log")
	c.AbortWithStatus(http.StatusInternalServerError)
}
