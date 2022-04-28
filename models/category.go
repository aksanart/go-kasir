package models

import (
	"aksan-kasir/config"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	CategoryId int       `gorm:"primaryKey;autoIncrement" json:"categoryId"`
	Name       string    `json:"name"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime:milli" json:"updatedAt"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type categoryList struct {
	CategoryId int    `gorm:"primaryKey" json:"categoryId,omitempty"`
	Name       string `json:"name,omitempty"`
}

type CategoryResult struct {
	Categories []categoryList `json:"Categories"`
	Meta       Meta           `json:"meta"`
}

// rename nama table
func (Category) TableName() string {
	return "category"
}

func FindAllCategory(db *gorm.DB, start int, limit int) (result CategoryResult) {
	var list []categoryList
	query := fmt.Sprintf("SELECT * FROM category LIMIT %v OFFSET %v", limit, start)
	rows, err := db.Raw(query).Rows()
	if err != nil {
		config.Throw(err.Error)
	}
	defer rows.Close()
	for rows.Next() {
		db.ScanRows(rows, &list)
	}
	result.Categories = list
	result.Meta.Limit = limit
	result.Meta.Skip = start
	result.Meta.Total = len(list)
	return result
}

func FindCategoryByID(db *gorm.DB, code string) (result categoryList) {
	query := fmt.Sprintf("SELECT * FROM category where category_id=%v", code)
	data := db.Raw(query).Scan(&result)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	return result
}

func FindCategoryByName(db *gorm.DB, code string) (result categoryList) {
	data := db.Raw("SELECT * FROM category where name=?", code).Scan(&result)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	return result
}

func CreateCategory(db *gorm.DB, category Category) Category {
	data := db.Create(&category)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	return category
}

func UpdateCategory(db *gorm.DB, category Category) {
	data := db.Model(category).Where("category_id = ?", category.CategoryId).Updates(category)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
}

func DeleteCategory(db *gorm.DB, data Category) {
	result := db.Model(data).Where("category_id = ?", data.CategoryId).Delete(data)
	if result.Error != nil {
		config.Throw(result.Error.Error)
	}
}
