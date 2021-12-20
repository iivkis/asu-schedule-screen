package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {

	gin.SetMode(gin.DebugMode)

	router := gin.Default()
	router.Use(gin.Logger(), gin.Recovery())

	router.GET("/screenshot", func(c *gin.Context) {
		photo, err := screenLink(c.Query("link"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		url, err := getUploadUrl()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		attachment, err := uploadAndSave(url, photo)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"attachment": attachment,
		})
	})

	router.Run()
}
