package models

import (
	"aksan-kasir/config"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type Cashier struct {
	CashierId int       `gorm:"primaryKey;autoIncrement" json:"cashierId"`
	Name      string    `json:"name"`
	Passcode  string    `json:"passcode"`
	Token     string    `json:"token"`
	Status    int       `json:"status"`
	UpdatedAt time.Time `gorm:"autoUpdateTime:milli" json:"updatedAt"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
}

type cashierList struct {
	CashierId int    `gorm:"primaryKey" json:"cashierId,omitempty"`
	Name      string `json:"name,omitempty"`
}

type loginList struct {
	Passcode string `json:"passcode,omitempty"`
	Token    string `json:"token,omitempty"`
}

type CashierResult struct {
	Cashiers []cashierList `json:"Cashiers"`
	Meta     Meta          `json:"meta"`
}

// rename nama table
func (Cashier) TableName() string {
	return "cashier"
}

func FindAllCashier(db *gorm.DB, start int, limit int) (result CashierResult) {
	var list []cashierList
	query := fmt.Sprintf("SELECT * FROM cashier LIMIT %v OFFSET %v", limit, start)
	rows, err := db.Raw(query).Rows()
	if err != nil {
		config.Throw(err.Error)
	}
	defer rows.Close()
	for rows.Next() {
		db.ScanRows(rows, &list)
	}
	result.Cashiers = list
	result.Meta.Limit = limit
	result.Meta.Skip = start
	result.Meta.Total = len(list)
	return result
}

func FindCashierByID(db *gorm.DB, code string) (result cashierList) {
	query := fmt.Sprintf("SELECT * FROM cashier where cashier_id=%v", code)
	data := db.Raw(query).Scan(&result)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	return result
}

func FindCashierByIDPasscode(db *gorm.DB, code string) (result loginList) {
	query := fmt.Sprintf("SELECT passcode FROM cashier where cashier_id=%v", code)
	data := db.Raw(query).Scan(&result)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	return result
}

func FindCashierByName(db *gorm.DB, code string) (result cashierList) {
	data := db.Raw("SELECT * FROM cashier where name=?", code).Scan(&result)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	return result
}

func CreateCashier(db *gorm.DB, cashier Cashier) Cashier {
	data := db.Create(&cashier)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	return cashier
}

func UpdateCashier(db *gorm.DB, cashier Cashier) {
	data := db.Model(cashier).Where("cashier_id = ?", cashier.CashierId).Updates(cashier)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
}

func DeleteCashier(db *gorm.DB, data Cashier) {
	result := db.Model(data).Where("cashier_id = ?", data.CashierId).Delete(data)
	if result.Error != nil {
		config.Throw(result.Error.Error)
	}
}

func Login(db *gorm.DB, code string, param config.CashierParamUpdate) (result loginList) {
	var cashier Cashier
	query := fmt.Sprintf("SELECT * FROM cashier where cashier_id=%v", code)
	data := db.Raw(query).Scan(&cashier)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	if cashier.Name == "" {
		config.Throw("not found")
	}
	if param.Passcode != cashier.Passcode {
		config.Throw("wrong passcode")
	}
	id := strconv.Itoa(cashier.CashierId)
	token, err := config.CreateToken(id, cashier.Name)
	if err != nil {
		config.Throw(err.Error)
	}
	cashier.Token = token
	cashier.Status = 1
	UpdateCashier(db, cashier)
	result.Token = cashier.Token
	return result
}

func Logout(db *gorm.DB, code string, param config.CashierParamUpdate) {
	var cashier Cashier
	query := fmt.Sprintf("SELECT * FROM cashier where cashier_id=%v", code)
	data := db.Raw(query).Scan(&cashier)
	if data.Error != nil {
		config.Throw(data.Error.Error)
	}
	if cashier.Name == "" {
		config.Throw("not found")
	}
	if param.Passcode != cashier.Passcode {
		config.Throw("wrong passcode")
	}
	cashier.Status = 2
	UpdateCashier(db, cashier)
}
