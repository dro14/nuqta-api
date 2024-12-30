package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
)

var errNoID = errors.New("id is not specified")

func failure(err error) gin.H {
	return gin.H{"error": err.Error()}
}
