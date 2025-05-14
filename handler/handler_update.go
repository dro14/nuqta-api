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
	broadcasters      = make(map[string][]chan *models.Event)
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
	channel := make(chan *models.Event)
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

	putMessages := make([]*models.Message, 0)
	deletedIds := make([]int64, 0)
	if c.Param("after") == "0" {
		now := time.Now().UnixMilli()

		chatUids, err := h.data.GetChats(ctx, uid, "private")
		if err != nil {
			sendSSEEvent(c, "error", err.Error())
			return
		}
		for _, chatUid := range chatUids {
			messages, err := h.data.GetMessages(ctx, uid, "private", chatUid, now)
			if err != nil {
				sendSSEEvent(c, "error", err.Error())
				return
			}
			putMessages = append(putMessages, messages...)
		}

		chatUids, err = h.data.GetChats(ctx, uid, "yordamchi")
		if err != nil {
			sendSSEEvent(c, "error", err.Error())
			return
		}
		for _, chatUid := range chatUids {
			messages, err := h.data.GetMessages(ctx, uid, "yordamchi", chatUid, now)
			if err != nil {
				sendSSEEvent(c, "error", err.Error())
				return
			}
			putMessages = append(putMessages, messages...)
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
		messages, err := h.data.GetUpdates(ctx, uid, "private", chatUids, after)
		if err != nil {
			sendSSEEvent(c, "error", err.Error())
			return
		}
		putMessages = append(putMessages, messages...)
		ids, err := h.data.GetDeletedMessages(ctx, uid, "private", chatUids, after)
		if err != nil {
			sendSSEEvent(c, "error", err.Error())
			return
		}
		deletedIds = append(deletedIds, ids...)

		chatUids, err = h.data.GetChats(ctx, uid, "yordamchi")
		if err != nil {
			sendSSEEvent(c, "error", err.Error())
			return
		}
		messages, err = h.data.GetUpdates(ctx, uid, "yordamchi", chatUids, after)
		if err != nil {
			sendSSEEvent(c, "error", err.Error())
			return
		}
		putMessages = append(putMessages, messages...)
		ids, err = h.data.GetDeletedMessages(ctx, uid, "yordamchi", chatUids, after)
		if err != nil {
			sendSSEEvent(c, "error", err.Error())
			return
		}
		for i, id := range ids {
			ids[i] = -1 * id
		}
		deletedIds = append(deletedIds, ids...)
	}

	if len(putMessages) > 0 {
		sendSSEEvent(c, "put_messages", putMessages)
	}

	if len(deletedIds) > 0 {
		sendSSEEvent(c, "delete_messages", deletedIds)
	}

	for {
		select {
		case event := <-channel:
			sendSSEEvent(c, event.Name, event.Data)
		case <-c.Request.Context().Done():
			return
		}
	}
}

func (h *Handler) ping(c *gin.Context) {
	broadcast(c.GetString("uid"), "pong", time.Now().UnixMilli())
}
