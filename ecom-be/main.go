package main

import (
    "ecom-be/routes"
    "github.com/gin-gonic/gin"
)


func main() {
	
    r := routes.SetupRouter()

    // Menjalankan server
    r.Run(":8080") // Menggunakan port 8080
}

