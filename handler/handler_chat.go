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
	messages, err := h.data.GetMessages(ctx, request.ChatUid, type_, request.Before)
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
	ctx := c.Request.Context()
	err = h.data.CreatePrivateMessage(ctx, message, c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	broadcast(message.RecipientUid, []*models.Message{message})
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
	err = h.data.ViewPrivateMessages(ctx, messages, c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	broadcast(messages[0].AuthorUid, messages)
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
	err = h.data.EditPrivateMessage(ctx, message, c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	broadcast(message.RecipientUid, []*models.Message{message})
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
	err = h.data.DeletePrivateMessage(ctx, message, c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
	broadcast(message.RecipientUid, []*models.Message{message})
	c.JSON(http.StatusOK, message)
}

func (h *Handler) removePrivateMessage(c *gin.Context) {
	message := &models.Message{}
	err := c.ShouldBindJSON(message)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	ctx := c.Request.Context()
	err = h.data.RemovePrivateMessage(ctx, message, c.GetString("uid"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}

func (h *Handler) typePrivateMessage(c *gin.Context) {
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
	err = h.data.CreateYordamchiMessage(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, request)

	go func() {
		conversation := make([]string, 0)
		for _, message := range messages {
			conversation = append(conversation, message.Text)
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, "firebase_uid", c.GetString("firebase_uid"))
		response, err := h.yordamchi.Respond(ctx, provider, conversation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
		response.ChatUid = request.ChatUid

		err = h.data.CreateYordamchiMessage(ctx, response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}

		broadcast(request.AuthorUid, []*models.Message{response})
	}()
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

	if len(messages) > 1 && len(messages) < 6 {
		c.JSON(http.StatusBadRequest, failure(e.ErrInvalidParams))
		return
	}

	ctx := c.Request.Context()
	request := messages[len(messages)-(1+len(messages)%2)]
	err = h.data.EditYordamchiMessage(ctx, request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, request)

	go func() {
		var id int64
		if len(messages)%2 != 0 {
			id = messages[len(messages)-1].Id
			messages = messages[:len(messages)-1]
		}

		conversation := make([]string, 0)
		for _, message := range messages {
			conversation = append(conversation, message.Text)
		}

		ctx := context.Background()
		ctx = context.WithValue(ctx, "firebase_uid", c.GetString("firebase_uid"))
		response, err := h.yordamchi.Respond(ctx, provider, conversation)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}
		response.Id = id

		err = h.data.EditYordamchiMessage(ctx, response)
		if err != nil {
			c.JSON(http.StatusInternalServerError, failure(err))
			return
		}

		broadcast(request.AuthorUid, []*models.Message{response})
	}()
}
