package controllers

import (
	"ecom-be/config"
	"ecom-be/middleware"
	"ecom-be/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateOrder membuat pesanan baru
func CreateOrder(c *gin.Context) {
	// Ambil user ID dari JWT token
	userClaims, _ := c.Get("user")
	claims := userClaims.(*middleware.Claims)
	userID := claims.UserID

	// Parse input
	var input struct {
		ShippingAddress string `json:"shipping_address" binding:"required"`
		PaymentMethod   string `json:"payment_method" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cari cart milik user
	var cart models.Cart
	if err := config.DB.Preload("CartItems.Product").Where("user_id = ?", userID).First(&cart).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Keranjang tidak ditemukan"})
		return
	}

	// Validasi: cart harus memiliki item
	if len(cart.CartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keranjang kosong"})
		return
	}

	// Hitung total
	var totalAmount float64 = 0
	for _, item := range cart.CartItems {
		totalAmount += item.Product.Price * float64(item.Quantity)
	}

	// Buat order baru
	order := models.Order{
		UserID:          userID,
		TotalAmount:     totalAmount,
		Status:          models.OrderStatusPending,
		ShippingAddress: input.ShippingAddress,
		PaymentMethod:   input.PaymentMethod,
	}

	// Transaction: create order & order items, update stock, clear cart
	tx := config.DB.Begin()

	// Buat order
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat pesanan"})
		return
	}

	// Buat order items & update stok
	for _, cartItem := range cart.CartItems {
		// Cek stok sekali lagi
		var product models.Product
		if err := tx.First(&product, cartItem.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
			return
		}

		if product.Stock < cartItem.Quantity {
			tx.Rollback()
			c.JSON(http.StatusBadRequest, gin.H{
				"error":     "Stok produk tidak mencukupi",
				"productId": product.ID,
			})
			return
		}

		// Buat order item
		orderItem := models.OrderItem{
			OrderID:   order.ID,
			ProductID: cartItem.ProductID,
			Quantity:  cartItem.Quantity,
			Price:     product.Price,
		}

		if err := tx.Create(&orderItem).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat item pesanan"})
			return
		}

		// Update stok produk
		product.Stock -= cartItem.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate stok produk"})
			return
		}
	}

	// Hapus semua cart item
	if err := tx.Where("cart_id = ?", cart.ID).Delete(&models.CartItem{}).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengosongkan keranjang"})
		return
	}

	// Commit transaksi
	tx.Commit()

	c.JSON(http.StatusCreated, gin.H{
		"message": "Pesanan berhasil dibuat",
		"order":   order,
	})
}

// GetOrders menampilkan semua pesanan user
func GetOrders(c *gin.Context) {
	// Ambil user ID dari JWT token
	userClaims, _ := c.Get("user")
	claims := userClaims.(*middleware.Claims)
	userID := claims.UserID

	// Parse query params untuk paginasi
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Query orders
	var orders []models.Order
	query := config.DB.Where("user_id = ?", userID).Order("created_at DESC").Offset(offset).Limit(limit)

	// Filter by status jika ada
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pesanan"})
		return
	}

	// Hitung total
	var total int64
	config.DB.Model(&models.Order{}).Where("user_id = ?", userID).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"meta": gin.H{
			"page":     page,
			"limit":    limit,
			"total":    total,
			"lastPage": (int(total) + limit - 1) / limit,
		},
	})
}

// GetOrderDetail menampilkan detail pesanan
func GetOrderDetail(c *gin.Context) {
	// Ambil user ID dari JWT token
	userClaims, _ := c.Get("user")
	claims := userClaims.(*middleware.Claims)
	userID := claims.UserID

	// Ambil order ID
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pesanan tidak valid"})
		return
	}

	// Query order dengan relasinya
	var order models.Order
	if err := config.DB.Preload("OrderItems.Product").Where("id = ? AND user_id = ?", orderId, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, order)
}

// CancelOrder membatalkan pesanan
func CancelOrder(c *gin.Context) {
	// Ambil user ID dari JWT token
	userClaims, _ := c.Get("user")
	claims := userClaims.(*middleware.Claims)
	userID := claims.UserID

	// Ambil order ID
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pesanan tidak valid"})
		return
	}

	// Cari order
	var order models.Order
	if err := config.DB.Where("id = ? AND user_id = ?", orderId, userID).First(&order).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	// Validasi: hanya bisa membatalkan pesanan dengan status pending atau processing
	if order.Status != models.OrderStatusPending && order.Status != models.OrderStatusProcessing {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Pesanan tidak dapat dibatalkan"})
		return
	}

	// Transaction: update order status & kembalikan stok
	tx := config.DB.Begin()

	// Load order items untuk mengembalikan stok
	var orderItems []models.OrderItem
	if err := tx.Where("order_id = ?", order.ID).Find(&orderItems).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data item pesanan"})
		return
	}

	// Kembalikan stok
	for _, item := range orderItems {
		var product models.Product
		if err := tx.First(&product, item.ProductID).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
			return
		}

		product.Stock += item.Quantity
		if err := tx.Save(&product).Error; err != nil {
			tx.Rollback()
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengembalikan stok produk"})
			return
		}
	}

	// Update status pesanan
	order.Status = models.OrderStatusCancelled
	if err := tx.Save(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membatalkan pesanan"})
		return
	}

	// Commit transaksi
	tx.Commit()

	c.JSON(http.StatusOK, gin.H{"message": "Pesanan berhasil dibatalkan"})
}

// GetAllOrders menampilkan semua pesanan (admin only)
func GetAllOrders(c *gin.Context) {
	// Parse query params untuk paginasi
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	// Query orders
	var orders []models.Order
	query := config.DB.Preload("User").Order("created_at DESC").Offset(offset).Limit(limit)

	// Filter by status jika ada
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pesanan"})
		return
	}

	// Hitung total
	var total int64
	config.DB.Model(&models.Order{}).Count(&total)

	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
		"meta": gin.H{
			"page":     page,
			"limit":    limit,
			"total":    total,
			"lastPage": (int(total) + limit - 1) / limit,
		},
	})
}

// UpdateOrderStatus mengubah status pesanan (admin only)
func UpdateOrderStatus(c *gin.Context) {
	// Ambil order ID
	orderId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pesanan tidak valid"})
		return
	}

	// Parse input
	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi status
	validStatus := map[string]bool{
		string(models.OrderStatusPending):    true,
		string(models.OrderStatusProcessing): true,
		string(models.OrderStatusShipped):    true,
		string(models.OrderStatusDelivered):  true,
		string(models.OrderStatusCancelled):  true,
	}

	if !validStatus[input.Status] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Status tidak valid"})
		return
	}

	// Cari order
	var order models.Order
	if err := config.DB.First(&order, orderId).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesanan tidak ditemukan"})
		return
	}

	// Update status
	order.Status = models.OrderStatus(input.Status)
	if err := config.DB.Save(&order).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate status pesanan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Status pesanan berhasil diubah",
		"order":   order,
	})
} 