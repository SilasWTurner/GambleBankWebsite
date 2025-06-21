package routes

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {
	// Define your routes here
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Welcome to GambleBank!",
		})
	})

	router.POST("/signup", Signup)
	router.POST("/login", Login)

	// WebSocket endpoint (not protected by AuthMiddleware)
	router.GET("/ws", WebSocketHandler)

	// Protected routes
	protected := router.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.POST("/send_invite", SendInvite)
		protected.GET("/list_invites", ListInvites)
		protected.POST("/accept_invite", AcceptInvite)
		protected.POST("/reject_invite", RejectInvite)
	}
}
