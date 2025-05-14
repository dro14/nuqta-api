package handler

import (
	"github.com/dro14/nuqta-service/models"
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
