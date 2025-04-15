package handler

import (
	"context"
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createResponse(c *gin.Context) {
	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	var request map[string]any
	err := c.ShouldBindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	list, ok := request["conversation"].([]any)
	if !ok || len(list) == 0 {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}
	var conversation []string
	for _, text := range list {
		conversation = append(conversation, text.(string))
	}
	language, ok := request["language"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}
	provider, ok := request["provider"].(string)
	if !ok {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	ctx = context.WithValue(ctx, "firebase_uid", firebaseUid)
	response, err := h.yordamchi.Respond(ctx, conversation, language, provider)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}
