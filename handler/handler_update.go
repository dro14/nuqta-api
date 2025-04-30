package handler

import (
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

var (
	broadcasters      = make(map[string][]chan any)
	broadcastersMutex sync.RWMutex
)

func (h *Handler) getUpdate(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")
	c.Writer.Header().Set("X-Accel-Buffering", "no")
	c.Writer.Flush()

	uid := c.GetString("uid")
	channel := make(chan any)
	broadcastersMutex.Lock()
	broadcasters[uid] = append(broadcasters[uid], channel)
	broadcastersMutex.Unlock()
	defer func() {
		broadcastersMutex.Lock()
		channels := broadcasters[uid]
		index := slices.Index(channels, channel)
		broadcasters[uid] = slices.Delete(channels, index, index+1)
		broadcastersMutex.Unlock()
		close(channel)
	}()

	after := c.Param("after")
	if after == "" || after == "0" {
		ctx := c.Request.Context()
		now := time.Now().UnixMilli()

		messages := make([]*models.Message, 0)
		chatUids, err := h.data.GetChats(ctx, uid, "private")
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
		for _, chatUid := range chatUids {
			chatMessages, err := h.data.GetMessages(ctx, chatUid, "private", now)
			if err != nil {
				c.JSON(http.StatusInternalServerError, failure(err))
				return
			}
			messages = append(messages, chatMessages...)
		}

		if len(messages) > 0 {
			sendSSEEvent(c, "messages", messages)
			messages = make([]*models.Message, 0)
		}

		chatUids, err = h.data.GetChats(ctx, uid, "yordamchi")
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
		for _, chatUid := range chatUids {
			chatMessages, err := h.data.GetMessages(ctx, chatUid, "yordamchi", now)
			if err != nil {
				c.JSON(http.StatusInternalServerError, failure(err))
				return
			}
			messages = append(messages, chatMessages...)
		}

		if len(messages) > 0 {
			sendSSEEvent(c, "messages", messages)
		}
	}

	for {
		select {
		case data := <-channel:
			switch data := data.(type) {
			case bool:
				sendSSEEvent(c, "pong", gin.H{"timestamp": time.Now().UnixMilli()})
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

func (h *Handler) ping(c *gin.Context) {
	broadcast(c.GetString("uid"), true)
	c.JSON(http.StatusOK, gin.H{"interval": 30})
}
