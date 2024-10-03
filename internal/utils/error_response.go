package utils

import "github.com/gin-gonic/gin"

func EncodeError(err error) gin.H {
	return gin.H{"error": err.Error()}
}
