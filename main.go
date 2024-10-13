package main

import (
	"choice/config"
	"choice/handler"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	config.InitMySQL()
	config.InitRedis()
	router := gin.Default()
	router.POST("/upload", handler.UploadExcel)
	router.GET("/records", handler.GetRecords)
	router.GET("/getrecordbyId/:id", handler.GetRecordByID)
	router.PUT("/edit/:id", handler.EditRecord)
	router.DELETE("/deleteById/:id", handler.DeleteRecordByID)
	router.DELETE("/deleteAll", handler.DeleteAllRecords)
	err := router.Run(":8080")
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)

	}
}
