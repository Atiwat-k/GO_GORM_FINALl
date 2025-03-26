package model

type CartItem struct {
	CartItemID int    `gorm:"primaryKey;autoIncrement" json:"cart_item_id"`
	CartID     int    `json:"cart_id"`
	ProductID  int    `json:"product_id"`
	Quantity   int    `json:"quantity"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Product    Product `gorm:"foreignKey:ProductID" json:"product"`
}
