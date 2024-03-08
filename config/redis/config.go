package redis

import "live-chat/config/config"

func getConfig() redisConfig {
	Info := redisConfig{
		Host:         "localhost",
		Port:         "6379",
		DB:           0,
		Password:     "",
		PoolSize:     30,
		MinIdleConns: 30,
	}
	if config.Config.IsSet("redis.host") {
		Info.Host = config.Config.GetString("redis.host")
	}
	if config.Config.IsSet("redis.port") {
		Info.Port = config.Config.GetString("redis.port")
	}
	if config.Config.IsSet("redis.db") {
		Info.DB = config.Config.GetInt("redis.db")
	}
	if config.Config.IsSet("redis.pass") {
		Info.Password = config.Config.GetString("redis.pass")
	}
	if config.Config.IsSet("redis.poolSize") {
		Info.PoolSize = config.Config.GetInt("redis.poolSize")
	}
	if config.Config.IsSet("redis.minIdleConns") {
		Info.MinIdleConns = config.Config.GetInt("redis.minIdleConns")
	}
	return Info
}
