package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) getUpdate(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Flush()

	ctx := c.Request.Context()
	uid := c.GetString("uid")
	chats, err := h.data.GetChats(ctx, uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	sendSSEEvent(c, "chats", chats)

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
