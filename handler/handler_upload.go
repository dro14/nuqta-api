package handler

import (
	"crypto/rand"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

func (h *Handler) upload(c *gin.Context) {
	filename := c.GetHeader("X-Filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, failure(fmt.Errorf("filename not provided in X-Filename header")))
		return
	}

	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	timestamp := time.Now().UnixNano()
	randomStr := make([]byte, 8)
	rand.Read(randomStr)
	extension := filepath.Ext(filename)
	filename = fmt.Sprintf("%d_%x%s", timestamp, randomStr, extension)

	err = os.WriteFile("uploads/"+filename, body, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.JSON(http.StatusOK, gin.H{"filename": filename})
}
