package controllers

import (
	"aksan-kasir/config"
	"aksan-kasir/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func ProductList(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			var param config.ProductParam
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			result := models.FindAllProduct(db, param)
			db.Commit()
			config.Response(c, "success", result)
		},
		Catch: func(e config.Exception) {
			db.Rollback()
			config.ErrorResponse(c, e)
		},
		Finally: func() {
			fmt.Println("selesai")
		},
	}.Do()
}

func ProductDetail(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("productId")
			if code == "" {
				config.Throw("ProductId is required")
			}
			result := models.FindProductByID(db, code)
			db.Commit()
			config.Response(c, "success", result)
		},
		Catch: func(e config.Exception) {
			db.Rollback()
			config.ErrorResponse(c, e)
		},
		Finally: func() {
			fmt.Println("selesai")
		},
	}.Do()
}

func ProductSave(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			var param config.ProductParamJson
			var param2 config.ProductParamJson
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			if param.Discount != param2.Discount {
				if param.Discount.ExpiredAt <= 0 || param.Discount.Qty <= 0 || param.Discount.Result <= 0 || param.Discount.Type == "" {
					config.Throw("discount required")
				}
			}
			result := models.CreateProduct(db, models.Product{CategoryId: param.CategoryID, Stock: param.Stock,
				Name: param.Name, Price: int64(param.Price), Image: param.Image,
				CreatedAt: time.Now(), UpdatedAt: time.Now()},
				models.Discount{Qty: param.Discount.Qty, Type: param.Discount.Type, Result: param.Discount.Result}, param.Discount.ExpiredAt)
			db.Commit()
			config.Response(c, "success", result)
		},
		Catch: func(e config.Exception) {
			db.Rollback()
			config.ErrorResponse(c, e)
		},
		Finally: func() {
			fmt.Println("selesai")
		},
	}.Do()
}

func ProductUpdate(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("productId")
			if code == "" {
				config.Throw("ProductId is required")
			}
			id, err := strconv.Atoi(code)
			if err != nil {
				config.Throw(err.Error)
			}
			var param config.ProductParamJsonUpdate
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			data := models.FindProductByID(db, code)
			if data.Name == "" {
				config.Throw("product tidak ditemukan")
			}
			models.UpdateProduct(db, param, id)
			db.Commit()
			config.Response(c, "success", nil)
		},
		Catch: func(e config.Exception) {
			db.Rollback()
			config.ErrorResponse(c, e)
		},
		Finally: func() {
			fmt.Println("selesai")
		},
	}.Do()
}

func ProductDelete(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("productId")
			if code == "" {
				config.Throw("ProductId is required")
			}
			codeInt, err := strconv.Atoi(code)
			if err != nil {
				config.Throw(err.Error)
			}
			result := models.FindProductByID(db, code)
			if result.Name == "" {
				config.Throw("not found")
			}
			models.DeleteProduct(db, models.Product{ProductId: codeInt})
			db.Commit()
			config.Response(c, "success", nil)
		},
		Catch: func(e config.Exception) {
			db.Rollback()
			config.ErrorResponse(c, e)
		},
		Finally: func() {
			fmt.Println("selesai")
		},
	}.Do()
}
