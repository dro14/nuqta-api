package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var (
	ErrNoUID = errors.New("uid is not specified")
)

func failure(err error) gin.H {
	return gin.H{"error": err.Error()}
}
