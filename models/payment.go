package models

import (
	"aksan-kasir/config"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type Payment struct {
	PaymentId int       `gorm:"primaryKey;autoIncrement" json:"paymentId"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Logo      string    `json:"logo"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updatedAt"`
}

type PaymentList struct {
	PaymentId int    `gorm:"primaryKey" json:"paymentId,omitempty"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Logo      string `json:"logo"`
	Card      []int  `json:"card,omitempty"`
}

type PaymentResult struct {
	Payments []PaymentList `json:"Payments"`
	Meta     Meta          `json:"meta"`
}

// rename nama table
func (Payment) TableName() string {
	return "payment"
}

func FindAllPayment(db *gorm.DB, param config.PaymentParam) (result PaymentResult) {
	var list []PaymentList
	query := fmt.Sprintf("SELECT * FROM payment LIMIT %v OFFSET %v", param.Limit, param.Skip)
	rows, err := db.Raw(query).Rows()
	if err != nil {
		config.Throw(err.Error)
	}
	defer rows.Close()
	for rows.Next() {
		db.ScanRows(rows, &list)
	}
	result.Payments = list
	result.Meta.Limit = param.Limit
	result.Meta.Skip = param.Skip
	result.Meta.Total = len(list)
	return result
}

func FindPaymentByID(db *gorm.DB, code string) (result PaymentList) {
	query := fmt.Sprintf("SELECT * FROM payment where payment_Id=%v", code)
	data := db.Raw(query).Scan(&result)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	return result
}

func FindPaymentByName(db *gorm.DB, code string) (result PaymentList) {
	data := db.Raw("SELECT * FROM payment where name=?", code).Scan(&result)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	return result
}

func CreatePayment(db *gorm.DB, payment Payment) Payment {
	data := db.Create(&payment)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	return payment
}

func UpdatePayment(db *gorm.DB, payment Payment) {
	data := db.Model(payment).Where("payment_id = ?", payment.PaymentId).Updates(payment)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
}

func DeletePayment(db *gorm.DB, data Payment) {
	result := db.Model(data).Where("payment_id = ?", data.PaymentId).Delete(data)
	if result.Error != nil {
		config.Throw(result.Error.Error)
	}
}
