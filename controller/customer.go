package controller

import (
	"go-grom/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var db *gorm.DB

// SetDB sets the database connection for the controller
func SetDB(database *gorm.DB) {
	db = database
}

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


