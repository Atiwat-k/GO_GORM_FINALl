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
	// รับข้อมูลจาก request body
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	// ดึง customerID จาก URL parameter
	customerIDStr := c.Param("customer_id")
	customerID, err := strconv.Atoi(customerIDStr) // แปลงจาก string เป็น int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid customer ID"})
		return
	}

	// ตรวจสอบว่า Cart ของลูกค้ามีอยู่ในฐานข้อมูลหรือไม่
	var cart model.Cart
	// ใช้ชื่อ table เป็น cart แทน carts
	if err := db.Where("customer_id = ?", customerID).First(&cart).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// ถ้าไม่พบ Cart ให้สร้างใหม่
			cart = model.Cart{
				CustomerID: customerID,
				CartName:   "My Cart", // สามารถกำหนดชื่อที่ต้องการ
			}
			// สร้าง Cart ใหม่ในฐานข้อมูล
			if err := db.Create(&cart).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error creating cart", "error": err.Error()})
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching cart", "error": err.Error()})
			return
		}
	}

	// ตรวจสอบการมีอยู่ของสินค้าในตาราง product
	var product model.Product
	if err := db.Where("product_id = ?", req.ProductID).First(&product).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Product not found", "error": err.Error()})
		return
	}

	// เพิ่มสินค้าไปยังรถเข็น
	var cartItem model.CartItem
	// ตรวจสอบว่ามีสินค้านี้ในรถเข็นแล้วหรือไม่
	if err := db.Where("cart_id = ? AND product_id = ?", cart.CartID, req.ProductID).First(&cartItem).Error; err == nil {
		// ถ้ามีแล้วให้เพิ่มจำนวนสินค้า
		cartItem.Quantity += req.Quantity
		if err := db.Save(&cartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating cart item", "error": err.Error()})
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
	// เพิ่มสินค้าใหม่ในฐานข้อมูล
	if err := db.Create(&newCartItem).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error adding item to cart", "error": err.Error()})
		return
	}

	// ตอบกลับว่าเพิ่มสินค้าเรียบร้อยแล้ว
	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
}
