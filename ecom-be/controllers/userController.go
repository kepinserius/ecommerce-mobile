package controllers

	import (
		"ecom-be/config"
		"ecom-be/models"
		"github.com/gin-gonic/gin"
		"golang.org/x/crypto/bcrypt"
		"github.com/golang-jwt/jwt/v4"
		"net/http"
		"time"
	)
	
	var jwtKey = []byte("secret_key_for_jwt") // Gantilah dengan key yang lebih aman
	
	// Struct untuk response login
	type LoginResponse struct {
		Token string `json:"token"`
	}
	
	// Fungsi untuk membuat JWT Token
	func GenerateJWT(user models.User) (string, error) {
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &jwt.RegisteredClaims{
			Issuer:    "ecomm-app",
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Subject:   string(user.ID),
		}
		
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		return token.SignedString(jwtKey)
	}
	
	// Login Handler untuk User/Admin
	func Login(c *gin.Context) {
		var input struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
	
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
	
		var user models.User
		if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
	
		// Verifikasi password
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
	
		// Generate JWT token
		token, err := GenerateJWT(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate token"})
			return
		}
	
		c.JSON(http.StatusOK, LoginResponse{Token: token})
	}
	
    "net/http"
    "golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
    var input models.User
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 10)
    input.Password = string(hashedPassword)

    config.DB.Create(&input)
    c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
}
