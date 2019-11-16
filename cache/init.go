package cache

import (
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"
	"vgo/model"
	"vgo/service"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

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
		redisClient = client
	}
}

// 存储结构体
func SetStruct(saveStruct interface{}, key string, expire time.Duration) {
	jsMa, _ := json.Marshal(saveStruct)
	redisClient.Set(key, jsMa, expire)
}

// 获取缓存的用户，如没有则从数据库中查询
// 获取一个用户信息
// 单个用户完整信息不存储在 session 中，故如果用户需要频繁使用该信息，则需要做缓存处理
// 缓存有效时间为1个小时，查询会先经过缓存，然后再向 Mysql 中查找
func GetUserFromRedis(key string) *model.User {
	cacheResult, _ := redisClient.Get(key).Result()

	user := &model.User{}
	if cacheResult == "" {
		userID, _ := strconv.Atoi(key)
		user.Id = userID
		model.DB.First(user)
		SetStruct(user, key, 300*time.Second)
	} else {
		json.Unmarshal([]byte(cacheResult), user)
	}
	return user
}

// 更新缓存信息
func UpdateUserCache(user *model.User) {
	usJs, _ := json.Marshal(user)
	ID := strconv.Itoa(user.Id)
	SetStruct(usJs, ID, 3600*time.Second)
}

// 删除缓存的用户信息
func DelCache(key string) {
	redisClient.Del(key)
}

// 获取注册服务
func GetRegisterUserFromRedis(key string) *service.UserRegisterService {
	result, _ := redisClient.Get(key).Result()
	register := &service.UserRegisterService{}
	json.Unmarshal([]byte(result), register)
	return register
}

// 普通的获取键值对
func Get(key string) string {
	result, _ := redisClient.Get(key).Result()
	return result
}

// 普通的存储键值对
func Set(key, value string, expire time.Duration) {
	redisClient.Set(key, value, expire)
}
