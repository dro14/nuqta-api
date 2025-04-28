package handler

import (
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

var channelsMap sync.Map

func (h *Handler) getUpdate(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.Flush()

	uid := c.GetString("uid")
	var channels []chan any
	value, ok := channelsMap.Load(uid)
	if ok {
		channels = value.([]chan any)
	}
	channel := make(chan any)
	channels = append(channels, channel)
	channelsMap.Store(uid, channels)
	defer func() {
		index := slices.Index(channels, channel)
		channels = slices.Delete(channels, index, index+1)
		channelsMap.Store(uid, channels)
		close(channel)
	}()

	messages := make([]*models.Message, 0)
	after := c.Param("after")
	if after == "" || after == "0" {
		ctx := c.Request.Context()
		chatUids, err := h.data.GetChats(ctx, uid, "private")
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
		now := time.Now().UnixMilli()
		for _, chatUid := range chatUids {
			chatMessages, err := h.data.GetMessages(ctx, chatUid, "private", now)
			if err != nil {
				c.JSON(http.StatusInternalServerError, failure(err))
				return
			}
			messages = append(messages, chatMessages...)
		}
	}

	sendSSEEvent(c, "messages", messages)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	i := 0
	for {
		select {
		case <-ticker.C:
			sendSSEEvent(c, "ping", gin.H{
				"ping":      i,
				"timestamp": time.Now().Add(5 * time.Hour).Format(time.DateTime),
			})
			i++
		case data := <-channel:
			switch data := data.(type) {
			case []*models.Message:
				sendSSEEvent(c, "messages", data)
			case string:
				sendSSEEvent(c, "typing", gin.H{"chat_uid": data})
			}
		case <-c.Request.Context().Done():
			return
		}
	}
}

func sendSSEEvent(c *gin.Context, event string, data any) {
	c.SSEvent(event, data)
	c.Writer.Flush()
}
