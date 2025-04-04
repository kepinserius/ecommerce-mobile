package models

import "time"

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusProcessing OrderStatus = "processing"
	OrderStatusShipped    OrderStatus = "shipped"
	OrderStatusDelivered  OrderStatus = "delivered"
	OrderStatusCancelled  OrderStatus = "cancelled"
)

type Order struct {
	ID            uint         `gorm:"primaryKey"`
	UserID        uint         `gorm:"not null"`
	User          User         `gorm:"foreignKey:UserID"`
	OrderItems    []OrderItem  `gorm:"foreignKey:OrderID"`
	TotalAmount   float64      `gorm:"not null"`
	Status        OrderStatus  `gorm:"type:varchar(20);default:'pending'"`
	ShippingAddress string     `gorm:"type:text;not null"`
	PaymentMethod string       `gorm:"type:varchar(50);not null"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type OrderItem struct {
	ID        uint    `gorm:"primaryKey"`
	OrderID   uint    `gorm:"not null"`
	ProductID uint    `gorm:"not null"`
	Product   Product `gorm:"foreignKey:ProductID"`
	Quantity  int     `gorm:"not null"`
	Price     float64 `gorm:"not null"`
} 