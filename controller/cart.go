package controller

import (
	"go-grom/model"
	"net/http"
	"strconv"
	"time"

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

	// ตรวจสอบสต็อกสินค้า
	if product.StockQuantity < req.Quantity {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not enough stock available"})
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
	} else {
		// ถ้าไม่พบให้เพิ่มสินค้าใหม่ลงไป
		newCartItem := model.CartItem{
			CartID:    cart.CartID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		// เพิ่มสินค้าใหม่ในฐานข้อมูล
		if err := db.Create(&newCartItem).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error adding item to cart", "error": err.Error()})
			return
		}
	}

	// ลดจำนวนสินค้าจากตาราง product
	product.StockQuantity -= req.Quantity
	if err := db.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating product stock", "error": err.Error()})
		return
	}

	// ตอบกลับว่าเพิ่มสินค้าเรียบร้อยแล้ว
	c.JSON(http.StatusOK, gin.H{"message": "Item added to cart successfully"})
}

// CartDetailResponse โครงสร้างข้อมูลที่ส่งกลับจาก API
type CartDetailResponse struct {
	CartID     int              `json:"cart_id"`
	CartName   string           `json:"cart_name"`
	CustomerID int              `json:"customer_id"`
	Items      []CartItemDetail `json:"items"`
}

// CartItemDetail โครงสร้างข้อมูลของสินค้าภายในรถเข็น
type CartItemDetail struct {
	ProductName string  `json:"product_name"`
	Quantity    int     `json:"quantity"`
	Price       float64 `json:"price"`
	TotalPrice  float64 `json:"total_price"`
}

// GetCartDetails ฟังก์ชันสำหรับดึงข้อมูลรถเข็นทั้งหมดของลูกค้า
func GetCartDetails(c *gin.Context) {
	// ดึง customerID จาก URL parameter
	customerIDStr := c.Param("customer_id")
	customerID, err := strconv.Atoi(customerIDStr) // แปลงจาก string เป็น int
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid customer ID"})
		return
	}

	// ดึงข้อมูล Cart ของลูกค้าทั้งหมด
	var carts []model.Cart
	if err := db.Where("customer_id = ?", customerID).Find(&carts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching carts", "error": err.Error()})
		return
	}

	var cartDetails []CartDetailResponse

	// วนลูปเพื่อดึงรายละเอียดของสินค้าในรถเข็นแต่ละคัน
	for _, cart := range carts {
		var cartDetail CartDetailResponse
		cartDetail.CartID = cart.CartID
		cartDetail.CartName = cart.CartName
		cartDetail.CustomerID = cart.CustomerID

		// ดึงรายการสินค้าจาก cart_item ที่เกี่ยวข้อง
		var cartItems []model.CartItem
		if err := db.Where("cart_id = ?", cart.CartID).Find(&cartItems).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching cart items", "error": err.Error()})
			return
		}

		// วนลูปเพื่อดึงข้อมูลของแต่ละรายการสินค้า
		for _, cartItem := range cartItems {
			var product model.Product
			if err := db.Where("product_id = ?", cartItem.ProductID).First(&product).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching product", "error": err.Error()})
				return
			}

			// คำนวณราคาของสินค้ารวม
			totalPrice := float64(cartItem.Quantity) * product.Price

			// เพิ่มรายละเอียดสินค้าไปยัง cartDetail
			cartItemDetail := CartItemDetail{
				ProductName: product.ProductName,
				Quantity:    cartItem.Quantity,
				Price:       product.Price,
				TotalPrice:  totalPrice,
			}

			cartDetail.Items = append(cartDetail.Items, cartItemDetail)
		}

		// เพิ่มรายละเอียดของ cart ลงในผลลัพธ์
		cartDetails = append(cartDetails, cartDetail)
	}

	// ส่งข้อมูลการตอบกลับ
	c.JSON(http.StatusOK, cartDetails)
}
