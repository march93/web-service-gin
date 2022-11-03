package main

import (
	"web-service-gin/database"
	"web-service-gin/server"
)

func main() {
	// Initialize database
	database.InitDB()

	// Setup routers
	router := server.NewRouter()
	router.Run("localhost:8080")
}
