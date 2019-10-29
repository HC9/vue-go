package cache

import (
	"log"
	"os"

	"github.com/go-redis/redis"
)

var RedisClient *redis.Client

// 初始化 redis 客户端
func init() {
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		log.Fatal(err)
	} else {
		RedisClient = client
	}
}
