package models

import "time"

type User struct {
    ID        uint      `gorm:"primaryKey"`
    Name      string    `gorm:"size:100;not null"`
    Email     string    `gorm:"unique;not null"`
    Password  string    `gorm:"not null"`
    Role      string    `gorm:"default:'user'"`
    CreatedAt time.Time
}
