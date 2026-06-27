package api

import (
	"github.com/gin-gonic/gin"
	"github.com/niljimeno/seamail/repository"
)

func hello(c *gin.Context) {
	mail, err := repository.ListMail()
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, mail)
}

func Run() {
	router := gin.Default()
	router.Use(CORSMiddleware())

	router.GET("/", hello)
	router.Run(":6000")
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
