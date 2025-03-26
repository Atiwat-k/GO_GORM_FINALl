package controller

import (
	"github.com/gin-gonic/gin"
)

// กำหนดฟังก์ชัน StartServer เพื่อเริ่มต้นเส้นทาง API
func StartServer(router *gin.Engine) {
	// เส้นทางแรก
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "API is now working",
		})
	})

	// เรียกใช้งาน DemoController

}

// DemoController จะเพิ่มเส้นทาง API สำหรับตัวอย่าง
