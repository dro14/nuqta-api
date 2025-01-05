package handler

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dro14/nuqta-service/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) upload(c *gin.Context) {
	filename := c.GetHeader("X-Filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoFilename))
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	randomStr := make([]byte, 8)
	rand.Read(randomStr)
	filename = fmt.Sprintf(
		"%d_%x%s",
		time.Now().UnixNano(),
		randomStr,
		filepath.Ext(filename),
	)

	err = os.WriteFile("storage/"+filename, body, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"filename": filename})
}
