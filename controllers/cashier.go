package controllers

import (
	"aksan-kasir/config"
	"aksan-kasir/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CashierList(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			var param config.CategoryParam
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			result := models.FindAllCashier(db, param.Skip, param.Limit)
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

func CashierDetail(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("cashierId")
			if code == "" {
				config.Throw("cashierId is required")
			}
			result := models.FindCashierByID(db, code)
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

func CashierSave(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			var param config.CashierParamJson
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			cek := models.FindCashierByName(db, param.Name)
			if cek.Name != "" {
				config.Throw("cashier dengan nama '" + param.Name + "' sudah ada")
			}
			result := models.CreateCashier(db, models.Cashier{Name: param.Name, Passcode: param.Passcode, CreatedAt: time.Now(), UpdatedAt: time.Now()})
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

func CashierUpdate(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("cashierId")
			if code == "" {
				config.Throw("cashierId is required")
			}
			var param config.CashierParamUpdate
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			result := models.FindCashierByID(db, code)
			if param.Name == "" {
				config.Throw("kategori tidak ditemukan")
			}
			models.UpdateCashier(db, models.Cashier{Name: param.Name, Passcode: param.Passcode, CashierId: result.CashierId, UpdatedAt: time.Now()})
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

func CashierDelete(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("cashierId")
			if code == "" {
				config.Throw("cashierId is required")
			}
			codeInt, err := strconv.Atoi(code)
			if err != nil {
				config.Throw(err.Error)
			}
			result := models.FindCashierByID(db, code)
			if result.Name == "" {
				config.Throw("not found")
			}
			models.DeleteCashier(db, models.Cashier{CashierId: codeInt})
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

func CashierPasscode(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("cashierId")
			if code == "" {
				config.Throw("cashierId is required")
			}
			result := models.FindCashierByIDPasscode(db, code)
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

func CashierLogin(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("cashierId")
			if code == "" {
				config.Throw("cashierId is required")
			}
			var param config.CashierParamUpdate
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			if param.Passcode == "" {
				config.Throw("passcode is required")
			}
			result := models.Login(db, code, param)
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

func CashierLogout(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("cashierId")
			if code == "" {
				config.Throw("cashierId is required")
			}
			var param config.CashierParamUpdate
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			if param.Passcode == "" {
				config.Throw("passcode is required")
			}
			models.Logout(db, code, param)
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
