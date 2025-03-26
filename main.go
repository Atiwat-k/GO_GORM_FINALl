package main

import (
	"fmt"
	"go-grom/controller"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	// Load config.yaml with viper
	viper.SetConfigName("config") // config file name (config.yaml)
	viper.AddConfigPath(".")      // config file path (same directory as Go code)
	viper.SetConfigType("yaml")   // config file format

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}

	// Get DSN from config.yaml
	dsn := viper.GetString("mysql.dsn")
	fmt.Println("Connecting to MySQL with DSN:", dsn)

	// Connect to database
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Test connection
	if err := db.Exec("SELECT 1").Error; err != nil {
		log.Fatalf("Database connection failed: %v", err)
	} else {
		fmt.Println("Connected to MySQL database successfully!")
	}

	// Set database for controller
	controller.SetDB(db)
}

func main() {
	// Initialize Gin
	r := gin.Default()

	// Start server with routes
	controller.StartServer(r)

	// Customer API routes
	r.GET("/customers", controller.GetCustomers)

	// Start server
	r.Run(":8080")
}