package handler

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) upload(c *gin.Context) {
	file, err := c.FormFile("upload")
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	timestamp := time.Now().UnixNano()
	randomStr := make([]byte, 8)
	rand.Read(randomStr)
	extension := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d_%x%s", timestamp, randomStr, extension)

	err = c.SaveUploadedFile(file, "uploads/"+filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"filename": filename})
}
