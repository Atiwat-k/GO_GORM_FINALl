package model

import "time"

// Product โครงสร้างข้อมูลสำหรับตาราง product
type Product struct {
	ProductID     int       `gorm:"primaryKey;autoIncrement" json:"product_id"`
	ProductName   string    `json:"product_name"`
	Description   string    `json:"description"`
	Price         float64   `json:"price"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName กำหนดชื่อของตารางในฐานข้อมูล
func (Product) TableName() string {
	return "product"
}
