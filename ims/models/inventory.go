package models

import (
	"time"
	"gorm.io/gorm"
)

type Inventory struct {
	ID uint`json:"id"`
	ProductID string`json:"product_id"`

	SKU string`json:"sku"`       
	Location  string`json:"location"`  

	Quantity  int`json:"quantity"`
	CreatedAt time.Time`json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
