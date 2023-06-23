package cache

import (
	"github.com/katerji/UserAuthKit/envs"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

var redisInstance *RedisClient

func GetRedisClient() *redis.Client {
	if redisInstance == nil {
		redisInstance = initRedis()
	}
	return redisInstance.client
}

func initRedis() *RedisClient {
	options, err := redis.ParseURL(envs.GetInstance().GetRedisURL())
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(options)
	redisInstance = &RedisClient{
		client: client,
	}
	return redisInstance
}
