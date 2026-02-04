package main

import (
	"log"
	"net/http"
	"os"

	"jpn-service/common"
	"jpn-service/routes"
)

func main() {
	// Load environment variables
	if err := common.LoadEnv(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Initialize database connection
	if err := common.InitDB(); err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer common.CloseDB()

	log.Println("✅ Connected to MongoDB Atlas successfully!")

	// Setup routes
	router := routes.SetupRoutes()

	// Get port from env or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start server
	log.Printf("🚀 Server is running on port %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
