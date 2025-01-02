package handler

import "github.com/gin-gonic/gin"

func failure(err error) gin.H {
	return gin.H{"error": err.Error()}
}
