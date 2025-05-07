package handler

import (
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

const version = "1.0.4"

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

	ctx := c.Request.Context()
	inviteCount, err := h.data.GetInviteCount(ctx, uid)
	if err != nil {
		sendSSEEvent(c, "error", err.Error())
		return
	}

	sendSSEEvent(c, "update", gin.H{
		"version":      version,
		"invite_count": inviteCount,
	})

	messages := make([]*models.Message, 0)
	if c.Param("after") == "0" {
		now := time.Now().UnixMilli()

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
	} else {
		after, err := strconv.ParseInt(c.Param("after"), 10, 64)
		if err != nil {
			sendSSEEvent(c, "error", err.Error())
			return
		}

		chatUids, err := h.data.GetChats(ctx, uid, "private")
		if err != nil {
			sendSSEEvent(c, "error", err.Error())
			return
		}
		chatMessages, err := h.data.GetUpdates(ctx, uid, "private", chatUids, after)
		if err != nil {
			sendSSEEvent(c, "error", err.Error())
			return
		}
		messages = append(messages, chatMessages...)

		chatUids, err = h.data.GetChats(ctx, uid, "yordamchi")
		if err != nil {
			sendSSEEvent(c, "error", err.Error())
			return
		}
		chatMessages, err = h.data.GetUpdates(ctx, uid, "yordamchi", chatUids, after)
		if err != nil {
			sendSSEEvent(c, "error", err.Error())
			return
		}
		messages = append(messages, chatMessages...)
	}

	if len(messages) > 0 {
		sendSSEEvent(c, "messages", messages)
	}

	for {
		select {
		case data := <-channel:
			switch data := data.(type) {
			case int64:
				sendSSEEvent(c, "pong", data)
			case []*models.Message:
				sendSSEEvent(c, "messages", data)
			case string:
				sendSSEEvent(c, "typing", data)
			}
		case <-c.Request.Context().Done():
			return
		}
	}
}

func (h *Handler) ping(c *gin.Context) {
	broadcast(c.GetString("uid"), time.Now().UnixMilli())
}
