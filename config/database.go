package config

import (
	"choice/models"
	"fmt"
	"log"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitMySQL() {
	user := viper.GetString("database.user")
	var err error
	DB, err = gorm.Open(mysql.Open(user), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to MySQL:", err)
	}
	fmt.Println("MySQL connected successfully")
	MigrateTable()
}

func MigrateTable() {
	DB.AutoMigrate(&models.Record{})
}
