package controllers

import (
	"ecom-be/config"
	"ecom-be/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// CreateProduct membuat produk baru (hanya admin)
func CreateProduct(c *gin.Context) {
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if result := config.DB.Create(&input); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat produk"})
		return
	}
	
	c.JSON(http.StatusCreated, gin.H{
		"message": "Produk berhasil dibuat",
		"product": input,
	})
}

// GetProducts menampilkan semua produk
func GetProducts(c *gin.Context) {
	var products []models.Product
	
	// Paginasi
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit
	
	// Query produk dengan paginasi
	query := config.DB.Offset(offset).Limit(limit)
	
	// Filter berdasarkan nama produk jika ada
	if search := c.Query("search"); search != "" {
		query = query.Where("name ILIKE ?", "%"+search+"%")
	}
	
	if err := query.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data produk"})
		return
	}
	
	// Hitung total produk
	var total int64
	config.DB.Model(&models.Product{}).Count(&total)
	
	c.JSON(http.StatusOK, gin.H{
		"products": products,
		"meta": gin.H{
			"page":      page,
			"limit":     limit,
			"total":     total,
			"lastPage":  (int(total) + limit - 1) / limit,
		},
	})
}

// GetProduct menampilkan produk berdasarkan ID
func GetProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")
	
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	
	c.JSON(http.StatusOK, product)
}

// UpdateProduct mengupdate produk berdasarkan ID (hanya admin)
func UpdateProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")
	
	// Cek apakah produk ada
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	
	// Bind input JSON ke product
	var input models.Product
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Update produk
	if err := config.DB.Model(&product).Updates(input).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengupdate produk"})
		return
	}
	
	// Ambil data produk yang telah diupdate
	config.DB.First(&product, id)
	
	c.JSON(http.StatusOK, gin.H{
		"message": "Produk berhasil diupdate",
		"product": product,
	})
}

// DeleteProduct menghapus produk berdasarkan ID (hanya admin)
func DeleteProduct(c *gin.Context) {
	var product models.Product
	id := c.Param("id")
	
	// Cek apakah produk ada
	if err := config.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produk tidak ditemukan"})
		return
	}
	
	// Hapus produk
	if err := config.DB.Delete(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menghapus produk"})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"message": "Produk berhasil dihapus"})
}
