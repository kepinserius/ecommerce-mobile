package routes

import (
	"ecom-be/controllers"
	"ecom-be/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowCredentials = true
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	r.Use(cors.New(config))

	// Rute tanpa autentikasi
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	// Rute untuk produk (publik)
	r.GET("/products", controllers.GetProducts)
	r.GET("/products/:id", controllers.GetProduct)

	// Rute dengan autentikasi
	authenticated := r.Group("/api")
	authenticated.Use(middleware.AuthRequired())
	{
		// Profil user
		authenticated.GET("/profile", controllers.GetProfile)

		// Cart
		authenticated.GET("/cart", controllers.GetCart)
		authenticated.POST("/cart", controllers.AddToCart)
		authenticated.PUT("/cart/:id", controllers.UpdateCartItem)
		authenticated.DELETE("/cart/:id", controllers.RemoveFromCart)
		authenticated.DELETE("/cart", controllers.ClearCart)

		// Pesanan
		authenticated.POST("/orders", controllers.CreateOrder)
		authenticated.GET("/orders", controllers.GetOrders)
		authenticated.GET("/orders/:id", controllers.GetOrderDetail)
		authenticated.PUT("/orders/:id/cancel", controllers.CancelOrder)
	}

	// Rute untuk admin
	admin := r.Group("/admin")
	admin.Use(middleware.AuthRequired())
	admin.Use(middleware.AdminRequired())
	{
		// Manajemen produk
		admin.POST("/products", controllers.CreateProduct)
		admin.PUT("/products/:id", controllers.UpdateProduct)
		admin.DELETE("/products/:id", controllers.DeleteProduct)

		// Manajemen pesanan
		admin.GET("/orders", controllers.GetAllOrders)
		admin.PUT("/orders/:id/status", controllers.UpdateOrderStatus)
	}

	return r
}
