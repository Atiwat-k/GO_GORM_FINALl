package controller

import (
	"go-grom/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddToCartRequest โครงสร้างข้อมูลที่รับจากผู้ใช้
type AddToCartRequest struct {
	ProductID int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" binding:"required"`
}

// AddToCart ฟังก์ชันสำหรับเพิ่มสินค้าไปยังรถเข็น
func AddToCart(c *gin.Context) {
	var req AddToCartRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	// ดึงข้อมูล customerID จาก parameter (สมมติว่ามาจาก URL parameter)
	customerIDStr := c.Param("customer_id")
	customerID, err := strconv.Atoi(customerIDStr) // แปลงจาก string เป็น int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid customer ID"})
		return
	}

	// ตรวจสอบว่ามี Cart ของลูกค้าหรือไม่
	var cart model.Cart
	if err := db.Where("customer_id = ?", customerID).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// ถ้าไม่พบ Cart ให้สร้างใหม่
			cart = model.Cart{
				CustomerID: customerID,
				CartName:   "My Cart", // สามารถกำหนดชื่อที่ต้องการ
			}
			if err := db.Create(&cart).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating cart"})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
			return
		}
	}

	// เพิ่มสินค้าไปยังรถเข็น
	var cartItem model.CartItem
	if err := db.Where("cart_id = ? AND product_id = ?", cart.CartID, req.ProductID).First(&cartItem).Error; err == nil {
		// ถ้ามีแล้วให้เพิ่มจำนวนสินค้า
		cartItem.Quantity += req.Quantity
		if err := db.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating cart item"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Cart item updated successfully"})
		return
	}

	// ถ้าไม่พบให้เพิ่มสินค้าใหม่ลงไป
	newCartItem := model.CartItem{
		CartID:    cart.CartID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
	}
	if err := db.Create(&newCartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error adding item to cart"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
}
