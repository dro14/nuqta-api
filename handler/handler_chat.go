package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

type MessageRequest struct {
	ChatUid  string   `json:"chat_uid"`
	ChatUids []string `json:"chat_uids"`
	Before   int64    `json:"before"`
	After    int64    `json:"after"`
}

func (h *Handler) createChat(c *gin.Context) {
	var request map[string]string
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if request["uid"] == "" || request["chat_with"] == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	chatUid, err := h.data.CreateChat(ctx, request["uid"], request["chat_with"])
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"chat_uid": chatUid})
}

func (h *Handler) getMessages(c *gin.Context) {
	type_ := c.Param("type")
	if type_ != "private" && type_ != "yordamchi" {
		c.JSON(http.StatusBadRequest, failure(e.ErrInvalidParams))
		return
	}

	request := &MessageRequest{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if request.ChatUid == "" || request.Before == 0 {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	messages, err := h.data.GetMessages(ctx, request.ChatUid, type_, request.Before)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *Handler) createPrivate(c *gin.Context) {
	message := &models.Message{}
	err := c.ShouldBindJSON(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}
	ctx := c.Request.Context()
	err = h.data.CreatePrivate(ctx, message, c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	broadcast(message.RecipientUid, []*models.Message{message})
	c.JSON(http.StatusOK, message)
}

func (h *Handler) viewPrivate(c *gin.Context) {
	var messages []*models.Message
	err := c.ShouldBindJSON(&messages)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}
	if len(messages) == 0 {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}
	ctx := c.Request.Context()
	err = h.data.ViewPrivate(ctx, messages, c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	broadcast(messages[0].AuthorUid, messages)
	c.JSON(http.StatusOK, messages)
}

func (h *Handler) editPrivate(c *gin.Context) {
	message := &models.Message{}
	err := c.ShouldBindJSON(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}
	ctx := c.Request.Context()
	err = h.data.EditPrivate(ctx, message, c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	broadcast(message.RecipientUid, []*models.Message{message})
	c.JSON(http.StatusOK, message)
}

func (h *Handler) deletePrivate(c *gin.Context) {
	message := &models.Message{}
	err := c.ShouldBindJSON(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}
	ctx := c.Request.Context()
	err = h.data.DeletePrivate(ctx, message, c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	broadcast(message.RecipientUid, []*models.Message{message})
	c.JSON(http.StatusOK, message)
}

func (h *Handler) removePrivate(c *gin.Context) {
	message := &models.Message{}
	err := c.ShouldBindJSON(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	ctx := c.Request.Context()
	err = h.data.RemovePrivate(ctx, message, c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}

func (h *Handler) typePrivate(c *gin.Context) {
	message := &models.Message{}
	err := c.ShouldBindJSON(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if message.ChatUid == "" || message.RecipientUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	broadcast(message.RecipientUid, message.ChatUid)
}

func (h *Handler) createYordamchi(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "openai" && provider != "google" {
		c.JSON(http.StatusBadRequest, failure(e.ErrInvalidParams))
		return
	}

	var messages []*models.Message
	err := c.ShouldBindJSON(&messages)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if len(messages) != 2 && len(messages) != 4 {
		c.JSON(http.StatusBadRequest, failure(e.ErrInvalidParams))
		return
	}

	ctx := c.Request.Context()
	request := messages[len(messages)-1]
	err = h.data.CreateYordamchi(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	go h.sendResponse(messages, c.GetString("firebase_uid"), provider)
	c.JSON(http.StatusOK, request)
}

func (h *Handler) editYordamchi(c *gin.Context) {
	provider := c.Param("provider")
	if provider != "openai" && provider != "google" {
		c.JSON(http.StatusBadRequest, failure(e.ErrInvalidParams))
		return
	}

	var messages []*models.Message
	err := c.ShouldBindJSON(&messages)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if len(messages) != 2 && len(messages) != 4 {
		c.JSON(http.StatusBadRequest, failure(e.ErrInvalidParams))
		return
	}

	ctx := c.Request.Context()
	request := messages[len(messages)-1]
	err = h.data.ClearYordamchi(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	err = h.data.CreateYordamchi(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	go h.sendResponse(messages, c.GetString("firebase_uid"), provider)
	c.JSON(http.StatusOK, request)
}

func (h *Handler) sendResponse(messages []*models.Message, firebaseUid, provider string) {
	request := messages[len(messages)-1]
	conversation := make([]string, 0)
	for _, message := range messages {
		conversation = append(conversation, message.Text)
	}

	now := time.Now()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "firebase_uid", firebaseUid)
	response, err := h.yordamchi.Respond(ctx, provider, conversation)
	if err != nil {
		response = &models.Message{
			Id:        now.UnixMicro(),
			ChatUid:   request.ChatUid,
			Timestamp: now.UnixMilli(),
			Text:      err.Error(),
		}
		broadcast(request.AuthorUid, []*models.Message{response})
		return
	}
	response.ChatUid = request.ChatUid

	err = h.data.CreateYordamchi(ctx, response)
	if err != nil {
		response = &models.Message{
			Id:        now.UnixMicro(),
			ChatUid:   request.ChatUid,
			Timestamp: now.UnixMilli(),
			Text:      err.Error(),
		}
		broadcast(request.AuthorUid, []*models.Message{response})
		return
	}

	broadcast(request.AuthorUid, []*models.Message{response})
}

func (h *Handler) deleteYordamchi(c *gin.Context) {
	message := &models.Message{}
	err := c.ShouldBindJSON(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}
	ctx := c.Request.Context()
	err = h.data.DeleteYordamchi(ctx, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}
