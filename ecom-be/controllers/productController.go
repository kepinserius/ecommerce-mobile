package controllers

import (
    "ecom-be/config"
    "ecom-be/models"

    "github.com/gin-gonic/gin"
    "net/http"
)

func CreateProduct(c *gin.Context) {
    var input models.Product
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    config.DB.Create(&input)
    c.JSON(http.StatusOK, gin.H{"message": "Product created"})
}

func GetProducts(c *gin.Context) {
    var products []models.Product
    if err := config.DB.Find(&products).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching products"})
        return
    }
    c.JSON(http.StatusOK, products)
}
