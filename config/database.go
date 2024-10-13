package config

import (
	"choice/models"
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMySQL() {
	dsn := "root:Nishanth123#@tcp(127.0.0.1:3306)/choiceDB?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}
	fmt.Println("MySQL connected successfully")
	MigrateTable()
}

func MigrateTable() {
	DB.AutoMigrate(&models.Record{})
}
