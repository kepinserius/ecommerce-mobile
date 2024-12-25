package routes

import (
	"ecomm-backend/controllers"
	"ecomm-backend/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Rute tanpa autentikasi
	r.POST("/login", controllers.Login) // Endpoint login

	// Rute dengan autentikasi menggunakan middleware
	// Middleware untuk memverifikasi JWT
	authenticated := r.Group("/api")
	authenticated.Use(middleware.AuthRequired())

	{
		// Produk
		authenticated.GET("/products", controllers.GetProducts)    // Menampilkan semua produk
		authenticated.GET("/products/:id", controllers.GetProduct) // Menampilkan produk berdasarkan ID
		authenticated.POST("/products", controllers.CreateProduct) // Membuat produk baru (hanya admin)
		authenticated.PUT("/products/:id", controllers.UpdateProduct) // Mengupdate produk (hanya admin)
		authenticated.DELETE("/products/:id", controllers.DeleteProduct) // Menghapus produk (hanya admin)
	}

	return r
}
