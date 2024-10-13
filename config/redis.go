package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var RDB *redis.Client

func InitRedis() {
	RDB = redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.address"),
		DB:   0,
	})
	c, err := RDB.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	fmt.Println("Redis Status", c)

}

func SetCache(key string, value string) error {
	return RDB.Set(context.Background(), key, value, 5*time.Minute).Err()
}

func GetCache(key string) (string, error) {
	return RDB.Get(context.Background(), key).Result()
}
