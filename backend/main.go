package main

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/silaswturner/GambleBankWebsite/backend/database"
	"github.com/silaswturner/GambleBankWebsite/backend/routes"
)

func main() {
	// Initialize the database and perform migrations
	database.Init()

	// Setup the router
	router := gin.Default()

	// Use the CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Update with your frontend's origin
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	routes.SetupRouter(router)

	// Start the backend server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
