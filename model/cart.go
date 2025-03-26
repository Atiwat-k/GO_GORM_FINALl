package model

type Cart struct {
	CartID     int       `gorm:"primaryKey;autoIncrement" json:"cart_id"`
	CustomerID int       `json:"customer_id"`
	CartName   string    `json:"cart_name"`
	CreatedAt  string    `json:"created_at"`
	UpdatedAt  string    `json:"updated_at"`
	Customer   Customer `gorm:"foreignKey:CustomerID" json:"customer"`
}
