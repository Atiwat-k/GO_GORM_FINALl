package controller

import (
	"fmt"
	"go-grom/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// GetCustomers retrieves all customers from the database
func GetCustomers(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Database not initialized"})
		return
	}

	var customers []model.Customer
	if err := db.Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error fetching customers"})
		return
	}
	c.JSON(http.StatusOK, customers)
}

// LoginRequest โครงสร้างข้อมูลที่รับจากผู้ใช้
type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login ฟังก์ชันตรวจสอบการเข้าสู่ระบบ
func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	fmt.Println("Login Attempt:", req.Email, req.Password) // Debug

	var customer model.Customer
	if err := db.Where("email = ?", req.Email).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		}
		return
	}

	fmt.Println("Stored Hash Password:", customer.Password)
	fmt.Println("Entered Password:", req.Password)

	// ตรวจสอบรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(req.Password)); err != nil {
		fmt.Println("Password mismatch error:", err)
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid email or password"})
		return
	}

	// ส่งข้อมูลผู้ใช้กลับ (ไม่รวมรหัสผ่าน)
	response := gin.H{
		"customer_id": customer.CustomerID,
		"first_name":  customer.FirstName,
		"last_name":   customer.LastName,
		"email":       customer.Email,
		"phone":       customer.PhoneNumber,
		"address":     customer.Address,
	}
	c.JSON(http.StatusOK, response)
}

// ChangePasswordRequest โครงสร้างข้อมูลที่รับจากผู้ใช้ในการเปลี่ยนรหัสผ่าน
type ChangePasswordRequest struct {
	OldPassword     string `json:"old_password" binding:"required"`
	NewPassword     string `json:"new_password" binding:"required"`
	ConfirmPassword string `json:"confirm_password" binding:"required"`
}

// ChangePassword ฟังก์ชันสำหรับการเปลี่ยนรหัสผ่าน
func ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid request"})
		return
	}

	// ตรวจสอบว่า รหัสผ่านใหม่และยืนยันรหัสผ่านตรงกันหรือไม่
	if req.NewPassword != req.ConfirmPassword {
		c.JSON(http.StatusBadRequest, gin.H{"message": "New password and confirmation password do not match"})
		return
	}

	// ดึงข้อมูลผู้ใช้จากฐานข้อมูล
	var customer model.Customer
	// สมมุติว่าคุณใช้ email เป็นตัวแปรระบุตัวผู้ใช้
	email := c.Param("email")
	if err := db.Where("email = ?", email).First(&customer).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Database error"})
		}
		return
	}

	// ตรวจสอบรหัสผ่านเก่ากับที่เก็บในฐานข้อมูล
	if err := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(req.OldPassword)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Incorrect old password"})
		return
	}

	// แฮชรหัสผ่านใหม่
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error hashing new password"})
		return
	}

	// อัพเดตรหัสผ่านในฐานข้อมูล
	customer.Password = string(hashedPassword)
	if err := db.Save(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Error updating password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Password updated successfully"})
}
