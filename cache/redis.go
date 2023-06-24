package cache

import (
	"github.com/katerji/UserAuthKit/envs"
	"github.com/redis/go-redis/v9"
)

type redisClient struct {
	client *redis.Client
}

var redisInstance *redisClient

func GetRedisClient() *redis.Client {
	if redisInstance == nil {
		redisInstance = initRedis()
	}
	return redisInstance.client
}

func initRedis() *redisClient {
	options, err := redis.ParseURL(envs.GetInstance().GetRedisURL())
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(options)
	redisInstance = &redisClient{
		client: client,
	}
	return redisInstance
}
