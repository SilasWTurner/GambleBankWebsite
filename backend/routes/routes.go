package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()

	// Define your routes here
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to GambleBank!",
		})
	})

	router.POST("/signup", Signup)
	router.POST("/login", Login)

	return router
}
