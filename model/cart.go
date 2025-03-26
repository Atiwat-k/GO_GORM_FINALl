package model

import (
	"time"
)

// Cart โครงสร้างข้อมูลของตาราง Cart
type Cart struct {
	CartID     int       `json:"cart_id" gorm:"primaryKey;autoIncrement"`
	CustomerID int       `json:"customer_id"`
	CartName   string    `json:"cart_name"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName กำหนดชื่อของตารางในฐานข้อมูล
func (Cart) TableName() string {
	return "cart"
}
