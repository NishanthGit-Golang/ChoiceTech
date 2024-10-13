package main

import (
	"choice/config"
	"choice/router"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	router.InitConfig()
	config.InitMySQL()
	config.InitRedis()
	r := gin.Default()
	router.SetupRouters(r)
	err := r.Run(viper.GetString("router.port"))
	if err != nil {
		log.Fatalf("Failed to start the server: %v", err)

	}
}
