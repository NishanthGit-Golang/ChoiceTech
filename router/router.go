package router

import (
	"choice/handler"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func SetupRouters(router *gin.Engine) {
	router.POST("/upload", handler.UploadExcel)
	router.GET("/records", handler.GetRecords)
	router.GET("/getrecordbyId/:id", handler.GetRecordByID)
	router.PUT("/edit/:id", handler.EditRecord)
	router.DELETE("/deleteById/:id", handler.DeleteRecordByID)
	router.DELETE("/deleteAll", handler.DeleteAllRecords)
	router.POST("/create", handler.CreateRecord)
}

func InitConfig() {
	viper.SetConfigFile("config.yaml") // Use a single YAML config file
	err := viper.ReadInConfig()        // Read the config file
	if err != nil {
		log.Fatalf("Error while reading config file: %s", err)
	}
}
