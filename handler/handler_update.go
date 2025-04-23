package handler

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getUpdate(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Flush()

	log.Printf("IP: %s\nUser-Agent: %s\n", c.ClientIP(), c.Request.UserAgent())

	sendSSEEvent(c, "connected", gin.H{
		"status":    "connected",
		"timestamp": time.Now().Unix(),
	})

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for i := 0; ; i++ {
		select {
		case <-ticker.C:
			sendSSEEvent(c, "update", gin.H{
				"update":    i,
				"timestamp": time.Now().Add(5 * time.Hour).Format(time.DateTime),
			})
		case <-c.Request.Context().Done():
			return
		}
	}
}

func sendSSEEvent(c *gin.Context, event string, data any) {
	c.SSEvent(event, data)
	c.Writer.Flush()
}
