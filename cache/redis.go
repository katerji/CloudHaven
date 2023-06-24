package cache

import (
	"context"
	"github.com/katerji/UserAuthKit/envs"
	"github.com/redis/go-redis/v9"
	"time"
)

type redisClient struct {
	redis *redis.Client
}

var redisInstance *redisClient

func GetRedisClient() *redisClient {
	if redisInstance == nil {
		redisInstance = initRedis()
	}
	return redisInstance
}

func initRedis() *redisClient {
	options, err := redis.ParseURL(envs.GetInstance().GetRedisURL())
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(options)
	redisInstance = &redisClient{
		redis: client,
	}
	return redisInstance
}

func (client *redisClient) Set(key string, value any, expiration time.Duration) error {
	return client.redis.Set(context.Background(), key, value, expiration).Err()
}

func (client *redisClient) Get(key string) string {
	return client.redis.Get(context.Background(), key).Val()
}

func (client *redisClient) Del(key string) bool {
	return client.redis.Del(context.Background(), key).Err() == nil
}

func (client *redisClient) TTL(key string) time.Duration {
	return client.redis.TTL(context.Background(), key).Val()
}

func (client *redisClient) GetMulti(keys []string) []any {
	return client.redis.MGet(context.Background(), keys...).Val()
}

func (client *redisClient) Keys(key string) []string {
	return client.redis.Keys(context.Background(), key).Val()
}
