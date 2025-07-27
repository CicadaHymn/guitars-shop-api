package api 

import (
	"github.com/gin-gonic/gin"
	"github.com/CicadaHymn/guitar-shop-api/internal/api/handlers"
)


func SetupRouters(r *gin.Engine) {
	r.POST("/order", handlers.CreateOrder)
	r.GET("/orders", handlers.GetOrders)
	
}

// Завтра нужно подключить SWAGGER для теста функций