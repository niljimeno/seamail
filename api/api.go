package api

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/niljimeno/seamail/repository"
)

func Run() {
	router := gin.Default()
	router.Use(CORSMiddleware())

	api := router.Group("/api")

	api.GET("/", listMail)
	api.GET("/mail/:id", readMail)
	api.DELETE("/mail/:id", removeMail)

	router.Run(":6000")
}

/* routes */

func listMail(c *gin.Context) {
	mail, err := repository.ListMail()
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, mail)
}

func readMail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, err)
		return
	}

	mail, err := repository.ReadMail(id)
	if err != nil {
		c.JSON(400, mail)
		return
	}

	c.JSON(200, mail)
}

func removeMail(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, err)
		return
	}

	err = repository.RemoveMail(id)
	if err != nil {
		c.JSON(400, err)
		return
	}

	c.JSON(200, "Removal succeeded")
}
