package models

import "time"

type Cart struct {
	ID        uint       `gorm:"primaryKey"`
	UserID    uint       `gorm:"not null;uniqueIndex"`
	User      User       `gorm:"foreignKey:UserID"`
	CartItems []CartItem `gorm:"foreignKey:CartID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CartItem struct {
	ID        uint    `gorm:"primaryKey"`
	CartID    uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int     `gorm:"not null;default:1"`
	CreatedAt time.Time
	UpdatedAt time.Time
} 