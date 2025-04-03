package handler

import (
	"net/http"

	"github.com/dro14/nuqta-service/e"
	"github.com/dro14/nuqta-service/models"
	"github.com/gin-gonic/gin"
)

func (h *Handler) createResponse(c *gin.Context) {
	firebaseUid := c.GetString("firebase_uid")
	if firebaseUid == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	request := &models.Request{}
	err := c.ShouldBindJSON(request)
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	if len(request.Conversation) == 0 || request.Language == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoParams))
		return
	}

	ctx := c.Request.Context()
	conversation := request.Conversation
	language := request.Language
	response, err := h.yordamchi.ProcessCompletions(ctx, conversation, language, firebaseUid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": response})
}
