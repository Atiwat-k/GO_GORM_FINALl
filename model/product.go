package model

type Product struct {
	ProductID     int     `gorm:"primaryKey;autoIncrement" json:"product_id"`
	ProductName   string  `json:"product_name"`
	Description   string  `json:"description"`
	Price         float64 `json:"price"`
	StockQuantity int     `json:"stock_quantity"`
	CreatedAt     string  `json:"created_at"`
	UpdatedAt     string  `json:"updated_at"`
}

func (Product) TableName() string {
	return "product"
}
