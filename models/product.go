package models

import (
	"aksan-kasir/config"
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ProductId  int `gorm:"primaryKey;autoIncrement" json:"productId"`
	CategoryId int
	DiscountID int
	SKU        string    `json:"sku"`
	Name       string    `json:"name"`
	Stock      int       `json:"stok"`
	Price      int64     `json:"price"`
	Image      string    `json:"image"`
	Category   Category  `gorm:"foreignkey:CategoryId" json:"category"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime:milli" json:"updatedAt"`
}

type Discount struct {
	DiscountID      int       `gorm:"primaryKey;autoIncrement" json:"discountId,omitempty"`
	Qty             int       `json:"qty,omitempty"`
	Type            string    `json:"type,omitempty"`
	Result          int       `json:"result,omitempty"`
	ExpiredAt       time.Time `json:"expiredAt,omitempty"`
	ExpiredAtFormat string    `json:"expiredAtFormat,omitempty"`
	StringFormat    string    `json:"stringFormat,omitempty"`
}

type DiscountList struct {
	DiscountID      int         `gorm:"primaryKey;autoIncrement" json:"discountId,omitempty"`
	Qty             int         `json:"qty,omitempty"`
	Type            string      `json:"type,omitempty"`
	Result          int         `json:"result,omitempty"`
	ExpiredAt       interface{} `json:"expiredAt,omitempty"`
	ExpiredAtFormat string      `json:"expiredAtFormat,omitempty"`
	StringFormat    string      `json:"stringFormat,omitempty"`
}

type ProductList struct {
	ProductId int          `gorm:"primaryKey" json:"productId"`
	SKU       string       `json:"sku"`
	Name      string       `json:"name"`
	Stock     int          `json:"stok"`
	Price     int64        `json:"price"`
	Image     string       `json:"image"`
	Category  categoryList `json:"category"`
	Discount  DiscountList `json:"discount"`
}

type Meta struct {
	Total int `json:"total"`
	Limit int `json:"limit"`
	Skip  int `json:"Skip"`
}

type ProductResult struct {
	Products []ProductList `json:"products"`
	Meta     Meta          `json:"meta"`
}

// rename nama table
func (Product) TableName() string {
	return "product"
}
func (Discount) TableName() string {
	return "discount"
}

func FindAllProduct(db *gorm.DB, param config.ProductParam) (result ProductResult) {
	var data []map[string]interface{}
	var rows []ProductList
	var res ProductList
	var limit, skip, q, catID string
	query := `SELECT product.product_id as p_id,product.sku as p_sku,product.name as p_name,
	product.stock as p_stock,product.price as p_price,product.image as p_image,
	category.category_id as c_id,category.name as c_name,discount.* 
	FROM product LEFT JOIN category ON product.category_id=category.category_id LEFT JOIN discount ON discount.discount_id=product.discount_id`
	if param.Q != "" {
		q = " product.name LIKE '%" + param.Q + "%'"
	}
	if param.CategoryID != 0 {
		catID = fmt.Sprintf(" product.category_id =%d", param.CategoryID)
	}
	if q != "" || catID != "" {
		and := catID + q
		if catID != "" && q != "" {
			and = catID + " AND " + q
		}
		query = query + " WHERE" + and
	}
	limit = fmt.Sprintf(" LIMIT %d", param.Limit)
	skip = fmt.Sprintf(" OFFSET %d", param.Skip)
	query = query + limit + skip
	err := db.Raw(query).Scan(&data).Error
	if err != nil {
		config.Throw(err.Error)
	}
	for _, v := range data {
		res.ProductId = int(v["p_id"].(int64))
		res.SKU = v["p_sku"].(string)
		res.Name = v["p_name"].(string)
		res.Stock = int(v["p_stock"].(int64))
		res.Price = v["p_price"].(int64)
		res.Image = v["p_image"].(string)
		if v["c_id"] != nil {
			res.Category.CategoryId = int(v["c_id"].(int64))
			res.Category.Name = v["c_name"].(string)
		}
		if v["discount_id"] != nil {
			res.Discount.DiscountID = int(v["discount_id"].(int64))
			res.Discount.Qty = int(v["qty"].(int64))
			res.Discount.Type = v["type"].(string)
			res.Discount.Result = int(v["result"].(int64))
			res.Discount.ExpiredAt = v["expired_at"]
			res.Discount.ExpiredAtFormat = v["expired_at_format"].(string)
			res.Discount.StringFormat = v["string_format"].(string)
		}
		rows = append(rows, res)
	}
	result.Products = rows
	result.Meta.Limit = param.Limit
	result.Meta.Skip = param.Skip
	result.Meta.Total = len(rows)
	return result
}

func FindProductByID(db *gorm.DB, code string) (result ProductList) {
	var res ProductList
	var list map[string]interface{}
	query := fmt.Sprintf(`SELECT product.product_id as p_id,product.sku as p_sku,product.name as p_name,
	product.stock as p_stock,product.price as p_price,product.image as p_image,
	category.category_id as c_id,category.name as c_name,discount.* 
	FROM product LEFT JOIN category ON product.category_id=category.category_id LEFT JOIN discount ON discount.discount_id=product.discount_id where product_id=%v`, code)
	data := db.Raw(query).Scan(&list)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	if list != nil {
		res.ProductId = int(list["p_id"].(int64))
		res.SKU = list["p_sku"].(string)
		res.Name = list["p_name"].(string)
		res.Stock = int(list["p_stock"].(int64))
		res.Price = list["p_price"].(int64)
		res.Image = list["p_image"].(string)
		if list["c_id"] != nil {
			res.Category.CategoryId = int(list["c_id"].(int64))
			res.Category.Name = list["c_name"].(string)
		}
		if list["discount_id"] != nil {
			res.Discount.DiscountID = int(list["discount_id"].(int64))
			res.Discount.Qty = int(list["qty"].(int64))
			res.Discount.Type = list["type"].(string)
			res.Discount.Result = int(list["result"].(int64))
			res.Discount.ExpiredAt = list["expired_at"]
			res.Discount.ExpiredAtFormat = list["expired_at_format"].(string)
			res.Discount.StringFormat = list["string_format"].(string)
		}
	}

	return res
}

// func FindProductByName(db *gorm.DB, code string) (result ProductList) {
// 	data := db.Raw("SELECT * FROM Product where name=?", code).Scan(&result)
// 	if data.Error != nil {
// 		config.Throw(data.Error.Error)
// 	}
// 	return result
// }

func CreateProduct(db *gorm.DB, product Product, diskon Discount, exp int) Product {
	if diskon.Qty > 0 {
		t, s := timeConverter(exp)
		diskon.ExpiredAtFormat = s
		diskon.ExpiredAt = t
		label := fmt.Sprintf("Buy %d only Rp. %d", diskon.Qty, diskon.Result)
		if diskon.Type == "PERCENT" {
			percenHarga := (diskon.Qty / 100) * int(product.Price)
			percen := product.Price - int64(percenHarga)
			label = fmt.Sprintf("Discount %d", diskon.Qty) + "%" + fmt.Sprintf(" Rp. %d", percen)
		}
		diskon.StringFormat = label
		err := db.Create(&diskon)
		if err.Error != nil {
			config.Throw(err.Error.Error)
		}
	}
	log.Println(diskon.DiscountID)
	product.DiscountID = diskon.DiscountID
	data := db.Create(&product)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	return product
}

func UpdateProduct(db *gorm.DB, param config.ProductParamJsonUpdate, id int) {
	var product Product
	if param.CategoryID > 0 {
		product.CategoryId = param.CategoryID
	}
	if param.Price > 0 {
		product.Price = int64(param.Price)
	}
	if param.Stock > 0 {
		product.Stock = param.Stock
	}
	if param.Image != "" {
		product.Image = param.Image
	}
	if param.Name != "" {
		product.Name = param.Name
	}
	data := db.Model(product).Where("product_id = ?", id).Updates(product)
	if data.Error != nil {
		log.Println(data.Error)
		config.Throw(data.Error.Error)
	}
}

func DeleteProduct(db *gorm.DB, data Product) {
	result := db.Model(data).Where("product_id = ?", data.ProductId).Delete(data)
	if result.Error != nil {
		config.Throw(result.Error.Error)
	}
}

func timeConverter(t int) (resTm time.Time, resStr string) {
	i, err := strconv.ParseInt("1405544146", 10, 64)
	if err != nil {
		panic(err)
	}
	tm := time.Unix(i, 0)
	str := tm.Format("02 January 2006")
	return tm, str
}
