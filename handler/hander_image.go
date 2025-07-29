package handler

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dro14/nuqta-service/utils/e"
	"github.com/gin-gonic/gin"
)

func (h *Handler) upload(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusBadRequest, failure(err))
		return
	}

	randomStr := make([]byte, 8)
	rand.Read(randomStr)
	ext := filepath.Ext(c.GetHeader("my-filename"))
	if ext == "" {
		ext = ".jpeg"
	}
	filename := fmt.Sprintf("%d_%x%s", time.Now().UnixNano(), randomStr, ext)

	err = os.WriteFile("images/"+filename, body, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}

	c.Header("my-filename", filename)
}

func (h *Handler) delete(c *gin.Context) {
	filename := c.GetHeader("my-filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, failure(e.ErrNoFilename))
		return
	}

	err := os.Remove("images/" + filename)
	if errors.Is(err, os.ErrNotExist) {
		log.Print("file not found: ", filename)
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, failure(err))
		return
	}
}
