package models

import (
	"aksan-kasir/config"
	"fmt"
	"strconv"
	"time"

	"github.com/jung-kurt/gofpdf"
	"gorm.io/gorm"
)

type Order struct {
	OrderID        int       `gorm:"primaryKey;autoIncrement" json:"orderId"`
	CashiersID     int       `json:"cashiersId"`
	PaymentTypesID int       `json:"paymentTypesId"`
	TotalPrice     int64     `json:"totalPrice"`
	TotalPaid      int64     `json:"totalPaid"`
	TotalReturn    int64     `json:"totalReturn"`
	ReceiptID      string    `gorm:"index:idx_name,unique" json:"receiptId"`
	IsDownload     bool      `json:"isDownload"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime:milli" json:"updatedAt"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type SubTotalList struct {
	ProductId        int          `json:"productId"`
	SKU              string       `json:"sku"`
	Name             string       `json:"name"`
	Stock            int          `json:"stok"`
	Price            int64        `json:"price"`
	Image            string       `json:"image"`
	ReceiptId        string       `json:"receiptId"`
	Discount         DiscountList `json:"diskon"`
	Qty              int          `json:"qty"`
	TotalNormalPrice int64        `json:"totalNormalPrice"`
	TotalFinalPrice  int64        `json:"totalFinalPrice"`
}

type OrderSaveList struct {
	ProductId        int          `json:"productId"`
	Name             string       `json:"name"`
	Price            int64        `json:"price"`
	Discount         DiscountList `json:"diskon"`
	Qty              int          `json:"qty"`
	TotalNormalPrice int64        `json:"totalNormalPrice"`
	TotalFinalPrice  int64        `json:"totalFinalPrice"`
}

type SubTotalResult struct {
	SubTotalList []SubTotalList `json:"products"`
	Subtotal     int64          `json:"subtotal"`
}

type OrderSaveResult struct {
	Order    Order          `json:"order"`
	Products []SubTotalList `json:"products"`
}

type OrderList struct {
	OrderID        int         `gorm:"primaryKey;autoIncrement" json:"orderId"`
	CashiersID     int         `json:"cashiersId"`
	PaymentTypesID int         `json:"paymentTypesId"`
	TotalPrice     int64       `json:"totalPrice"`
	TotalPaid      int64       `json:"totalPaid"`
	TotalReturn    int64       `json:"totalReturn"`
	ReceiptID      string      `gorm:"index:idx_name,unique" json:"receiptId"`
	UpdatedAt      time.Time   `gorm:"autoUpdateTime:milli" json:"updatedAt"`
	CreatedAt      time.Time   `gorm:"autoCreateTime" json:"createdAt"`
	Cashier        cashierList `json:"cashier"`
	PaymentType    PaymentList `json:"payment_type"`
}
type OrderListResult struct {
	Orders []OrderList `json:"orders"`
	Meta   Meta        `json:"meta"`
}

// rename nama table
func (Order) TableName() string {
	return "order"
}

func FindSubTotal(db *gorm.DB, param []config.OrderParam) (result SubTotalResult) {
	var res SubTotalList
	var arr []SubTotalList
	var total, last int64
	for _, v := range param {
		var list map[string]interface{}
		query := fmt.Sprintf(`SELECT product.product_id as p_id,product.sku as p_sku,product.name as p_name,
	product.stock as p_stock,product.price as p_price,product.image as p_image,
	category.category_id as c_id,category.name as c_name,discount.* 
	FROM product LEFT JOIN category ON product.category_id=category.category_id LEFT JOIN discount ON discount.discount_id=product.discount_id 
	WHERE product.product_id=%d`, v.ProductId)
		data := db.Raw(query).Scan(&list)
		if data.Error != nil {
			config.Throw(data.Error.Error)
		}
		if list != nil {
			total = int64(v.Qty) * res.Price
			res.ProductId = int(list["p_id"].(int64))
			res.SKU = list["p_sku"].(string)
			res.Name = list["p_name"].(string)
			res.Stock = int(list["p_stock"].(int64))
			res.Price = list["p_price"].(int64)
			res.Image = list["p_image"].(string)
			if list["discount_id"] != nil {
				res.Discount.DiscountID = int(list["discount_id"].(int64))
				res.Discount.Qty = int(list["qty"].(int64))
				res.Discount.Type = list["type"].(string)
				res.Discount.Result = int(list["result"].(int64))
				res.Discount.ExpiredAt = list["expired_at"]
				res.Discount.ExpiredAtFormat = list["expired_at_format"].(string)
				res.Discount.StringFormat = list["string_format"].(string)
				if res.Discount.Qty == v.Qty {
					percenHarga := (res.Discount.Qty / 100) * int(res.Price)
					total = res.Price - int64(percenHarga)
				}
			}
			last = last + total
			res.Qty = v.Qty
			res.TotalNormalPrice = int64(v.Qty) * res.Price
			res.TotalFinalPrice = total
			arr = append(arr, res)
		}
	}
	result.Subtotal = last
	result.SubTotalList = arr
	return result
}

func OrderSave(db *gorm.DB, param config.OrderParamJson) (result OrderSaveResult) {
	var order Order
	var res SubTotalList
	var arr []SubTotalList
	var total, last int64
	for _, v := range param.Products {
		var list map[string]interface{}
		query := fmt.Sprintf(`SELECT product.product_id as p_id,product.sku as p_sku,product.name as p_name,
	product.stock as p_stock,product.price as p_price,product.image as p_image,
	category.category_id as c_id,category.name as c_name,discount.* 
	FROM product LEFT JOIN category ON product.category_id=category.category_id LEFT JOIN discount ON discount.discount_id=product.discount_id 
	WHERE product.product_id=%d`, v.ProductId)
		data := db.Raw(query).Scan(&list)
		if data.Error != nil {
			config.Throw(data.Error.Error)
		}
		if list != nil {
			total = int64(v.Qty) * res.Price
			res.ProductId = int(list["p_id"].(int64))
			res.SKU = list["p_sku"].(string)
			res.Name = list["p_name"].(string)
			res.Stock = int(list["p_stock"].(int64))
			res.Price = list["p_price"].(int64)
			res.Image = list["p_image"].(string)
			if list["discount_id"] != nil {
				res.Discount.DiscountID = int(list["discount_id"].(int64))
				res.Discount.Qty = int(list["qty"].(int64))
				res.Discount.Type = list["type"].(string)
				res.Discount.Result = int(list["result"].(int64))
				res.Discount.ExpiredAt = list["expired_at"]
				res.Discount.ExpiredAtFormat = list["expired_at_format"].(string)
				res.Discount.StringFormat = list["string_format"].(string)
				if res.Discount.Qty == v.Qty {
					percenHarga := (res.Discount.Qty / 100) * int(res.Price)
					total = res.Price - int64(percenHarga)
				}
			}
			last = last + total
			res.Qty = v.Qty
			res.TotalNormalPrice = int64(v.Qty) * res.Price
			res.TotalFinalPrice = total
			arr = append(arr, res)
			stok := res.Stock - v.Qty
			UpdateProduct(db, config.ProductParamJsonUpdate{
				Stock: stok,
			}, res.ProductId)
		}
	}
	kembalian := param.TotalPaid - last
	if kembalian < 0 {
		config.Throw("uang tidak cukup")
	}
	order.CashiersID = param.CashierId
	order.PaymentTypesID = param.PaymentId
	order.TotalPrice = last
	order.TotalPaid = param.TotalPaid
	order.TotalReturn = kembalian
	order.ReceiptID = config.GenerateID(6)
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	CreateOrder(db, order)
	result.Order = order
	result.Products = arr
	return result
}

func CreateOrder(db *gorm.DB, table Order) Order {
	data := db.Create(&table)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	return table
}

func FindOrderList(db *gorm.DB, param config.CategoryParam) (result OrderListResult) {
	var list []OrderList
	var data []map[string]interface{}
	order := "`order`"
	query := fmt.Sprintf(`SELECT o.*,cashier.name AS c_name,cashier.cashier_id AS c_id,
	 payment.payment_id AS p_id, payment.name AS p_name, payment.logo AS p_logo, payment.type AS p_type
	 FROM `+order+` As o 
	 LEFT JOIN cashier ON cashier.cashier_id=o.cashiers_id 
	 LEFT JOIN payment ON payment.payment_id=o.payment_types_id
	 LIMIT %v OFFSET %v`, param.Limit, param.Skip)
	err := db.Raw(query).Scan(&data).Error
	if err != nil {
		config.Throw(err.Error)
	}
	for _, v := range data {
		var res OrderList
		res.OrderID = int(v["order_id"].(int64))
		res.CashiersID = int(v["cashiers_id"].(int64))
		res.PaymentTypesID = int(v["payment_types_id"].(int64))
		res.TotalPrice = v["total_price"].(int64)
		res.TotalPaid = v["total_paid"].(int64)
		res.TotalReturn = v["total_return"].(int64)
		res.ReceiptID = v["receipt_id"].(string)
		res.CreatedAt = v["created_at"].(time.Time)
		res.UpdatedAt = v["updated_at"].(time.Time)
		res.Cashier.CashierId = int(v["c_id"].(int64))
		res.Cashier.Name = v["c_name"].(string)
		res.PaymentType.PaymentId = int(v["p_id"].(int64))
		res.PaymentType.Name = v["p_name"].(string)
		res.PaymentType.Logo = v["p_logo"].(string)
		res.PaymentType.Type = v["p_type"].(string)
		list = append(list, res)
	}
	result.Orders = list
	result.Meta.Limit = int(param.Limit)
	result.Meta.Skip = int(param.Skip)
	result.Meta.Total = len(list)
	return result
}

func FindOrderDetail(db *gorm.DB, id int) (result OrderList) {
	var data map[string]interface{}
	order := "`order`"
	query := fmt.Sprintf(`SELECT o.*,cashier.name AS c_name,cashier.cashier_id AS c_id,
	 payment.payment_id AS p_id, payment.name AS p_name, payment.logo AS p_logo, payment.type AS p_type
	 FROM `+order+` As o 
	 LEFT JOIN cashier ON cashier.cashier_id=o.cashiers_id 
	 LEFT JOIN payment ON payment.payment_id=o.payment_types_id
	 WHERE o.order_id=%v`, id)
	err := db.Raw(query).Scan(&data).Error
	if err != nil {
		config.Throw(err.Error)
	}
	var res OrderList
	if data == nil {
		config.Throw("Order Not Found")
	}
	res.OrderID = int(data["order_id"].(int64))
	res.CashiersID = int(data["cashiers_id"].(int64))
	res.PaymentTypesID = int(data["payment_types_id"].(int64))
	res.TotalPrice = data["total_price"].(int64)
	res.TotalPaid = data["total_paid"].(int64)
	res.TotalReturn = data["total_return"].(int64)
	res.ReceiptID = data["receipt_id"].(string)
	res.CreatedAt = data["created_at"].(time.Time)
	res.UpdatedAt = data["updated_at"].(time.Time)
	res.Cashier.CashierId = int(data["c_id"].(int64))
	res.Cashier.Name = data["c_name"].(string)
	res.PaymentType.PaymentId = int(data["p_id"].(int64))
	res.PaymentType.Name = data["p_name"].(string)
	res.PaymentType.Logo = data["p_logo"].(string)
	res.PaymentType.Type = data["p_type"].(string)
	return res
}

func OrderPdf(db *gorm.DB, id int) error {
	data := FindOrderDetail(db, id)
	strID := strconv.Itoa(id)
	strTotalPrice := strconv.Itoa(int(data.TotalPrice))
	strTotalPaid := strconv.Itoa(int(data.TotalPaid))
	strTotalReturn := strconv.Itoa(int(data.TotalReturn))
	nameFile := "Order_" + strID + ".pdf"
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Text(40, 10, "Order id: "+strID)
	pdf.Text(40, 30, "Total price: "+strTotalPrice)
	pdf.Text(40, 50, "Total paid: "+strTotalPaid)
	pdf.Text(40, 70, "Total return: "+strTotalReturn)

	err := pdf.OutputFileAndClose("./" + nameFile)
	if err != nil {
		config.Throw(err.Error())
	}
	UpdateOrder(db, Order{OrderID: id, IsDownload: true})
	return nil
}

func UpdateOrder(db *gorm.DB, order Order) {
	data := db.Model(order).Where("order_id = ?", order.OrderID).Updates(order)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
}

func FindOrderByID(db *gorm.DB, code string) (result Order) {
	query := fmt.Sprintf("SELECT * FROM `order` AS o where order_id=%v", code)
	data := db.Raw(query).Scan(&result)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	var result2 Order
	if result == result2 {
		config.Throw("Order Not Found")
	}
	return result
}
