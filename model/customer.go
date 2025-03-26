package model

import (
	"time"
)

// Customer เป็นโมเดลสำหรับตาราง customer ในฐานข้อมูล
type Customer struct {
	CustomerID  int       `gorm:"column:customer_id;AUTO_INCREMENT;primary_key"`
	FirstName   string    `gorm:"column:first_name;NOT NULL"`
	LastName    string    `gorm:"column:last_name;NOT NULL"`
	Email       string    `gorm:"column:email;unique;NOT NULL"`
	PhoneNumber string    `gorm:"column:phone_number"`
	Address     string    `gorm:"column:address"`
	Password    string    `gorm:"column:password;NOT NULL"`
	CreatedAt   time.Time `gorm:"column:created_at;default:CURRENT_TIMESTAMP"`
	UpdatedAt   time.Time `gorm:"column:updated_at;default:CURRENT_TIMESTAMP"`
}

// TableName กำหนดชื่อของตารางที่ใช้ในฐานข้อมูล
func (m *Customer) TableName() string {
	return "customer"
}
