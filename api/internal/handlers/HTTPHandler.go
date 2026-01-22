package handlers

import (
	"github.com/gin-gonic/gin"
)

func OptionsCORSHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Status(204)
	}
}
