package main

import (
	"aksan-kasir/config"
	"aksan-kasir/controllers"
	"aksan-kasir/models"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	log.Fatal(err.Error())
	// }
	config.InitDB()
	migrate()
	r := gin.Default()
	r.Use(gin.Logger())
	r.GET("categories", controllers.CategoryList)
	r.GET("categories/:categoryId", controllers.CategoryDetail)
	r.POST("categories", controllers.CategorySave)
	r.PUT("categories/:categoryId", controllers.CategoryUpdate)
	r.DELETE("categories/:categoryId", controllers.CategoryDelete)

	r.GET("products", controllers.ProductList)
	r.GET("products/:productId", controllers.ProductDetail)
	r.POST("products", controllers.ProductSave)
	r.PUT("products/:productId", controllers.ProductUpdate)
	r.DELETE("products/:productId", controllers.ProductDelete)

	r.POST("cashiers", controllers.CashierSave)
	r.PUT("cashiers/:cashierId", controllers.CashierUpdate)
	r.DELETE("cashiers/:cashierId", controllers.CashierDelete)
	r.GET("cashiers", controllers.CashierList)
	r.GET("cashiers/:cashierId", controllers.CashierDetail)

	r.GET("cashiers/:cashierId/passcode", controllers.CashierPasscode)
	r.POST("cashiers/:cashierId/login", controllers.CashierLogin)
	r.POST("cashiers/:cashierId/logout", controllers.CashierLogout)

	r.POST("payments", controllers.PaymentSave)
	r.PUT("payments/:paymentId", controllers.PaymentUpdate)
	r.DELETE("payments/:paymentId", controllers.PaymentDelete)
	r.GET("payments", controllers.PaymentList)
	r.GET("payments/:paymentId", controllers.PaymentDetail)

	r.POST("orders/subtotal", controllers.OrderSubTotal)
	r.POST("orders", controllers.OrderSave)
	r.GET("orders", controllers.OrderList)
	r.GET("orders/:orderId", controllers.OrderDetail)
	r.GET("orders/:orderId/download", controllers.OrderDownload)
	r.GET("orders/:orderId/check-download", controllers.OrderCekDownload)

	r.Run(":3030")
}

func migrate() {
	for _, model := range models.MigrateModels() {
		err := config.DB.Debug().AutoMigrate(model.Model)

		if err != nil {
			log.Fatal(err)
		}
	}
}
