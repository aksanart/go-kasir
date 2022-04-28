package config

import (
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() (database *gorm.DB, err error) {
	var dsn strings.Builder
	dsn.WriteString(os.Getenv("MYSQL_USER") + ":" + os.Getenv("MYSQL_PASSWORD") + "@tcp(" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT") + ")" + "/" + os.Getenv("MYSQL_DBNAME") + "?charset=utf8mb4&parseTime=True&loc=Local")
	database, err = gorm.Open(mysql.Open(dsn.String()), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return database, nil
}

func InitDB() {
	conn, err := ConnectDB()
	if err != nil {
		log.Fatal(err.Error())
	}
	DB = conn
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal(err.Error())
	}
	sqlDB.SetMaxIdleConns(100)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(5 * time.Second)
	DB.Debug().AutoMigrate()
	log.Println("connected to :" + os.Getenv("MYSQL_HOST") + ":" + os.Getenv("MYSQL_PORT"))
}
