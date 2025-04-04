package controllers

import (
	"ecom-be/config"
	"ecom-be/middleware"
	"ecom-be/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCart mengambil cart user yang sedang login
func GetCart(c *gin.Context) {
	// Ambil user ID dari JWT token
	userClaims, _ := c.Get("user")
	claims := userClaims.(*middleware.Claims)
	userID := claims.UserID

	// Cek apakah user memiliki cart
	var cart models.Cart
	result := config.DB.Preload("CartItems.Product").Where("user_id = ?", userID).First(&cart)
	
	// Jika cart tidak ditemukan, buat cart baru
	if result.Error != nil {
		cart = models.Cart{
			UserID: userID,
		}
		config.DB.Create(&cart)
	}

	// Hitung total harga cart
	var total float64 = 0
	for _, item := range cart.CartItems {
		total += item.Product.Price * float64(item.Quantity)
	}

	c.JSON(http.StatusOK, gin.H{
		"cart": cart,
		"total": total,
	})
}

// AddToCart menambahkan produk ke keranjang
func AddToCart(c *gin.Context) {
	// Ambil user ID dari JWT token
	userClaims, _ := c.Get("user")
	claims := userClaims.(*middleware.Claims)
	userID := claims.UserID

	// Parse input
	var input struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek ketersediaan produk
	var product models.Product
	if err := config.DB.First(&product, input.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	// Cek stok produk
	if product.Stock < input.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stok produk tidak mencukupi"})
		return
	}

	// Cari atau buat cart untuk user
	var cart models.Cart
	result := config.DB.Where("user_id = ?", userID).First(&cart)
	if result.Error != nil {
		cart = models.Cart{UserID: userID}
		config.DB.Create(&cart)
	}

	// Cek apakah produk sudah ada di cart
	var cartItem models.CartItem
	result = config.DB.Where("cart_id = ? AND product_id = ?", cart.ID, input.ProductID).First(&cartItem)
	
	if result.Error != nil {
		// Produk belum ada di cart, buat cart item baru
		cartItem = models.CartItem{
			CartID:    cart.ID,
			ProductID: input.ProductID,
			Quantity:  input.Quantity,
		}
		config.DB.Create(&cartItem)
	} else {
		// Produk sudah ada di cart, update quantity
		cartItem.Quantity += input.Quantity
		config.DB.Save(&cartItem)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil ditambahkan ke keranjang"})
}

// UpdateCartItem mengubah jumlah produk di keranjang
func UpdateCartItem(c *gin.Context) {
	// Ambil user ID dari JWT token
	userClaims, _ := c.Get("user")
	claims := userClaims.(*middleware.Claims)
	userID := claims.UserID

	// Parse input
	var input struct {
		Quantity int `json:"quantity" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ambil ID cart item
	cartItemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID item tidak valid"})
		return
	}

	// Cari cart milik user
	var cart models.Cart
	if err := config.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Keranjang tidak ditemukan"})
		return
	}

	// Cari cart item
	var cartItem models.CartItem
	if err := config.DB.Where("id = ? AND cart_id = ?", cartItemID, cart.ID).First(&cartItem).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item tidak ditemukan di keranjang"})
		return
	}

	// Cek ketersediaan stok
	var product models.Product
	if err := config.DB.First(&product, cartItem.ProductID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}

	if product.Stock < input.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Stok produk tidak mencukupi"})
		return
	}

	// Update quantity
	cartItem.Quantity = input.Quantity
	config.DB.Save(&cartItem)

	c.JSON(http.StatusOK, gin.H{"message": "Jumlah produk berhasil diubah"})
}

// RemoveFromCart menghapus produk dari keranjang
func RemoveFromCart(c *gin.Context) {
	// Ambil user ID dari JWT token
	userClaims, _ := c.Get("user")
	claims := userClaims.(*middleware.Claims)
	userID := claims.UserID

	// Ambil ID cart item
	cartItemID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID item tidak valid"})
		return
	}

	// Cari cart milik user
	var cart models.Cart
	if err := config.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Keranjang tidak ditemukan"})
		return
	}

	// Hapus cart item
	result := config.DB.Where("id = ? AND cart_id = ?", cartItemID, cart.ID).Delete(&models.CartItem{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Item tidak ditemukan di keranjang"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus dari keranjang"})
}

// ClearCart menghapus semua produk dari keranjang
func ClearCart(c *gin.Context) {
	// Ambil user ID dari JWT token
	userClaims, _ := c.Get("user")
	claims := userClaims.(*middleware.Claims)
	userID := claims.UserID

	// Cari cart milik user
	var cart models.Cart
	if err := config.DB.Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Keranjang tidak ditemukan"})
		return
	}

	// Hapus semua cart item
	config.DB.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})

	c.JSON(http.StatusOK, gin.H{"message": "Keranjang berhasil dikosongkan"})
} 