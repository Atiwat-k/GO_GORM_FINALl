package model

import (
	"time"
)

// CartItem โครงสร้างข้อมูลของตาราง CartItem
type CartItem struct {
	CartItemID int       `json:"cart_item_id" gorm:"primaryKey;autoIncrement"`
	CartID     int       `json:"cart_id"`
	ProductID  int       `json:"product_id"`
	Quantity   int       `json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// TableName กำหนดชื่อของตารางในฐานข้อมูล
func (CartItem) TableName() string {
	return "cart_item"
}
