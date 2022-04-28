package controllers

import (
	"aksan-kasir/config"
	"aksan-kasir/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func PaymentList(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			var param config.PaymentParam
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			result := models.FindAllPayment(db, param)
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

func PaymentDetail(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("paymentId")
			if code == "" {
				config.Throw("paymentId is required")
			}
			result := models.FindPaymentByID(db, code)
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

func PaymentSave(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			var param config.PaymentParamJson
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			cek := models.FindPaymentByName(db, param.Name)
			if cek.Name != "" {
				config.Throw("payment dengan nama '" + param.Name + "' sudah ada")
			}
			result := models.CreatePayment(db, models.Payment{Name: param.Name, Type: param.Type, Logo: param.Logo, CreatedAt: time.Now(), UpdatedAt: time.Now()})
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

func PaymentUpdate(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("paymentId")
			if code == "" {
				config.Throw("paymentId is required")
			}
			var param config.PaymentParamUpdate
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			result := models.FindPaymentByID(db, code)
			if param.Name == "" {
				config.Throw("kategori tidak ditemukan")
			}
			models.UpdatePayment(db, models.Payment{Name: param.Name, Type: param.Type, Logo: param.Logo, PaymentId: result.PaymentId, UpdatedAt: time.Now()})
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

func PaymentDelete(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("paymentId")
			if code == "" {
				config.Throw("paymentId is required")
			}
			codeInt, err := strconv.Atoi(code)
			if err != nil {
				config.Throw(err.Error)
			}
			result := models.FindPaymentByID(db, code)
			if result.Name == "" {
				config.Throw("not found")
			}
			models.DeletePayment(db, models.Payment{PaymentId: codeInt})
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
