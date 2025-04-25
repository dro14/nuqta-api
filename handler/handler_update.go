package handler

import (
	"net/http"
	"time"

	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) getUpdate(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Flush()

	messages := make([]*models.Message, 0)
	after := c.Param("after")
	if after == "" || after == "0" {
		ctx := c.Request.Context()
		uid := c.GetString("uid")
		chatUids, err := h.data.GetChats(ctx, uid, "private")
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
		now := time.Now().UnixMilli()
		for _, chatUid := range chatUids {
			chatMessages, err := h.data.GetMessages(ctx, "private", chatUid, now)
			if err != nil {
				c.JSON(http.StatusInternalServerError, failure(err))
				return
			}
			messages = append(messages, chatMessages...)
		}
	}

	sendSSEEvent(c, "messages", messages)

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
