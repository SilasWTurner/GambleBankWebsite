package main

import (
	"log"

	"github.com/silaswturner/GambleBankWebsite/backend/database"
	"github.com/silaswturner/GambleBankWebsite/backend/routes"
)

func main() {
	// Initialize the database and perform migrations
	database.Init()

	// Setup the router
	router := routes.SetupRouter()

	// Start the backend server
	if err := router.Run(":8080"); err != nil {
		log.Fatalf("Failed to start the server: %v", err)
	}
}
