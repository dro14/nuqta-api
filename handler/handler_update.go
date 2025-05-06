package handler

import (
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
			sendSSEEvent(c, "error", err.Error())
			return
		}
		for _, chatUid := range chatUids {
			chatMessages, err := h.data.GetMessages(ctx, uid, "private", chatUid, now)
			if err != nil {
				sendSSEEvent(c, "error", err.Error())
				return
			}
			messages = append(messages, chatMessages...)
		}

		chatUids, err = h.data.GetChats(ctx, uid, "yordamchi")
		if err != nil {
			sendSSEEvent(c, "error", err.Error())
			return
		}
		for _, chatUid := range chatUids {
			chatMessages, err := h.data.GetMessages(ctx, uid, "yordamchi", chatUid, now)
			if err != nil {
				sendSSEEvent(c, "error", err.Error())
				return
			}
			messages = append(messages, chatMessages...)
		}

		if len(messages) > 0 {
			sendSSEEvent(c, "messages", messages)
		}

		inviteCount, err := h.data.GetInviteCount(ctx, uid)
		if err != nil {
			sendSSEEvent(c, "error", err.Error())
			return
		}
		sendSSEEvent(c, "invite_count", gin.H{"count": inviteCount})
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
}
