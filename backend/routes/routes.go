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

	// Protected routes
	protected := router.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.GET("/play", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Welcome to the play area!",
			})
		})

		protected.GET("/leaderboard", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Welcome to the leaderboard!",
			})
		})
	}

	return router
}
