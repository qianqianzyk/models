package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

type redisConfig struct {
	Host         string `yaml:"host"`
	Port         string `yaml:"port"`
	DB           int    `yaml:"db"`
	User         string `yaml:"user"`
	Password     string `yaml:"pass"`
	PoolSize     int    `yaml:"poolSize"`
	MinIdleConns int    `yaml:"minIdleConns"`
}

var RedisClient *redis.Client
var RedisInfo redisConfig

func init() {
	info := getConfig()

	RedisClient = redis.NewClient(&redis.Options{
		Addr:         info.Host + ":" + info.Port,
		Password:     info.Password,
		DB:           info.DB,
		PoolSize:     info.PoolSize,
		MinIdleConns: info.MinIdleConns,
	})

	RedisInfo = info

	// 测试连接
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
}
