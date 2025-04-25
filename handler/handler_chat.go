package handler

import (
	"context"
	"net/http"

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

func (h *Handler) getMessageList(c *gin.Context) {
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
	messages, err := h.data.GetMessages(ctx, type_, request.ChatUid, request.Before)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, messages)
}

func (h *Handler) createPrivateMessage(c *gin.Context) {
	message := &models.Message{}
	err := c.ShouldBindJSON(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}
	channel, ok := channels.Load(message.RecipientUid)
	if ok {
		channel.(chan models.Message) <- *message
	}
	ctx := c.Request.Context()
	err = h.data.CreateMessage(ctx, message, "private")
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	c.JSON(http.StatusOK, message)
}

func (h *Handler) viewPrivateMessages(c *gin.Context) {
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
	err = h.data.ViewPrivateMessages(ctx, messages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	c.JSON(http.StatusOK, messages)
}

func (h *Handler) editPrivateMessage(c *gin.Context) {
	message := &models.Message{}
	err := c.ShouldBindJSON(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}
	ctx := c.Request.Context()
	err = h.data.EditPrivateMessage(ctx, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	c.JSON(http.StatusOK, message)
}

func (h *Handler) deletePrivateMessage(c *gin.Context) {
	message := &models.Message{}
	err := c.ShouldBindJSON(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}
	ctx := c.Request.Context()
	err = h.data.DeletePrivateMessage(ctx, message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	c.JSON(http.StatusOK, message)
}

func (h *Handler) createYordamchiMessage(c *gin.Context) {
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
	err = h.data.CreateMessage(ctx, request, "yordamchi")
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	conversation := make([]string, len(messages))
	for i, message := range messages {
		conversation[i] = message.Text
	}

	ctx = context.WithValue(ctx, "firebase_uid", c.GetString("firebase_uid"))
	response, err := h.yordamchi.Respond(ctx, provider, conversation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	response.ChatUid = request.ChatUid
	response.InReplyTo = request.Id

	err = h.data.CreateMessage(ctx, response, "yordamchi")
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, []*models.Message{request, response})
}

func (h *Handler) editYordamchiMessage(c *gin.Context) {
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

	if len(messages) != 3 && len(messages) != 5 {
		c.JSON(http.StatusBadRequest, failure(e.ErrInvalidParams))
		return
	}

	ctx := c.Request.Context()
	request := messages[len(messages)-2]
	err = h.data.EditYordamchiMessage(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	conversation := make([]string, len(messages))
	for i, message := range messages[:len(messages)-1] {
		conversation[i] = message.Text
	}

	ctx = context.WithValue(ctx, "firebase_uid", c.GetString("firebase_uid"))
	response, err := h.yordamchi.Respond(ctx, provider, conversation)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	messages[len(messages)-1].Text = response.Text
	messages[len(messages)-1].AuthorUid = response.AuthorUid
	response = messages[len(messages)-1]

	err = h.data.EditYordamchiMessage(ctx, response)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, []*models.Message{request, response})
}
