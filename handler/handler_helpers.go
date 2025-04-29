package handler

import (
	"github.com/gin-gonic/gin"
)

func failure(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func sendSSEEvent(c *gin.Context, event string, data any) {
	c.SSEvent(event, data)
	c.Writer.Flush()
}

func broadcast(uid string, data any) {
	broadcastersMutex.RLock()
	channels := broadcasters[uid]
	for _, channel := range channels {
		channel <- data
	}
	broadcastersMutex.RUnlock()
}
