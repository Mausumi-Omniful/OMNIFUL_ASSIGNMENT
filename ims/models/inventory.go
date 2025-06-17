package models

import (
	"time"
	"gorm.io/gorm"
)

type Inventory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	ProductID   string         `gorm:"not null" json:"product_id"`
	SKU         string         `gorm:"unique;not null" json:"sku"`
	Quantity    int            `gorm:"not null" json:"quantity"`
	Location    string         `json:"location"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
