package main

import (
	"ecom-be/config"
	"ecom-be/routes"
	"fmt"
	"log"
	"os"
)

func main() {
	// Load environment variables
	config.LoadEnv()
	
	// Connect ke database
	config.ConnectDatabase()
	
	// Setup router
	r := routes.SetupRouter()
	
	// Set port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	// Jalankan server
	log.Printf("Server berjalan di port %s", port)
	if err := r.Run(":" + port); err != nil {
		fmt.Printf("Gagal menjalankan server: %v", err)
	}
}

