package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/niljimeno/seamail/config"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		pass := c.GetHeader("AUTH")

		if pass != config.Password {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Next()
	}
}
