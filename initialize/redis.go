package initialize

import (
	"context"
	"github.com/biningo/boil-gin/global"
	"github.com/go-redis/redis/v8"
)

/**
*@Author lyer
*@Date 4/19/21 10:09
*@Describe
**/

func InitRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     global.G_CONFIG.Redis.Addr,
		DB:       global.G_CONFIG.Redis.DB,
		Password: global.G_CONFIG.Redis.Password,
		Username: global.G_CONFIG.Redis.Username,
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		panic(err)
	}
	return client
}
