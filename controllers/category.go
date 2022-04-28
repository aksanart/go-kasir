package controllers

import (
	"aksan-kasir/config"
	"aksan-kasir/models"
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func CategoryList(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			var param config.CategoryParam
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			result := models.FindAllCategory(db, param.Skip, param.Limit)
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

func CategoryDetail(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("categoryId")
			if code == "" {
				config.Throw("categoryId is required")
			}
			result := models.FindCategoryByID(db, code)
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

func CategorySave(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			var param config.CategoryParamJson
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			cek := models.FindCategoryByName(db, param.Name)
			if cek.Name != "" {
				config.Throw("kategori dengan nama '" + param.Name + "' sudah ada")
			}
			result := models.CreateCategory(db, models.Category{Name: param.Name, CreatedAt: time.Now(), UpdatedAt: time.Now()})
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

func CategoryUpdate(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("categoryId")
			if code == "" {
				config.Throw("categoryId is required")
			}
			var param config.CategoryParamJson
			if err := c.ShouldBind(&param); err != nil {
				config.Throw(err.Error())
			}
			result := models.FindCategoryByID(db, code)
			if param.Name == "" {
				config.Throw("kategori tidak ditemukan")
			}
			models.UpdateCategory(db, models.Category{Name: param.Name, CategoryId: result.CategoryId, UpdatedAt: time.Now()})
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

func CategoryDelete(c *gin.Context) {
	db := config.DB.Begin()
	config.Block{
		Try: func() {
			code := c.Param("categoryId")
			if code == "" {
				config.Throw("categoryId is required")
			}
			codeInt, err := strconv.Atoi(code)
			if err != nil {
				config.Throw(err.Error)
			}
			result := models.FindCategoryByID(db, code)
			if result.Name == "" {
				config.Throw("not found")
			}
			models.DeleteCategory(db, models.Category{CategoryId: codeInt})
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
