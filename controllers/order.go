package controllers

import (
	"aksan-kasir/config"
	"aksan-kasir/models"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func OrderSubTotal(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			var param []config.OrderParam
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			result := models.FindSubTotal(db, param)
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

func OrderSave(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			var param config.OrderParamJson
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			auth := c.GetHeader("Authorization")
			if auth == "" {
				config.Throw("authorization is required")
			}
			tokenString := strings.Replace(auth, "Bearer ", "", -1)
			token, err := config.ExtractTokenMetadata(tokenString)
			if err != nil {
				config.Throw(err.Error())
			}
			param.CashierId = int(token.CashierId)
			result := models.OrderSave(db, param)
			db.Commit()
			config.Response(c, "success", result)
		},
		Catch: func(e config.Exception) {
			db.Rollback()
			log.Println(e)
			config.ErrorResponse(c, e)
		},
		Finally: func() {
			fmt.Println("selesai")
		},
	}.Do()
}

func OrderList(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			var param config.CategoryParam
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			result := models.FindOrderList(db, param)
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

func OrderDetail(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("orderId")
			if code == "" {
				config.Throw("orderId is required")
			}
			id, err := strconv.Atoi(code)
			if err != nil {
				config.Throw(err.Error)
			}
			result := models.FindOrderDetail(db, id)
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

func OrderDownload(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("orderId")
			if code == "" {
				config.Throw("orderId is required")
			}
			id, err := strconv.Atoi(code)
			if err != nil {
				config.Throw(err.Error)
			}
			err = models.OrderPdf(db, id)
			if err != nil {
				config.Throw(err.Error)
			}
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

func OrderCekDownload(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("orderId")
			if code == "" {
				config.Throw("orderId is required")
			}
			result := models.FindOrderByID(db, code)
			db.Commit()
			config.Response(c, "success", struct {
				IsDownload bool `json:"isDownload"`
			}{
				result.IsDownload,
			})
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
