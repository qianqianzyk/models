package configSecret

import (
	"context"
	"fmt"
	"live-chat/app/models"
	"live-chat/config/database"
	"live-chat/config/redis"
	"time"
)

var ctx = context.Background()

func getConfig(key string) (string, error) {
	val, err := redis.RedisClient.Get(ctx, key).Result()
	if err == nil {
		return val, nil
	}
	fmt.Println(err)
	var config models.Config
	result := database.DB.Model(&models.Config{}).Where(
		"title = ?", key).First(&config)
	if result.Error != nil {
		return "111", result.Error
	}
	redis.RedisClient.Set(ctx, key, config.Content, 0)
	return config.Content, nil
}

func setConfig(key, value string) error {
	redis.RedisClient.Set(ctx, key, value, 0)
	res := database.DB.Model(models.Config{}).Where(
		&models.Config{
			Title: key,
		}).Updates(&models.Config{
		Title:   key,
		Content: value,
	})
	if res.RowsAffected == 0 {
		rc := database.DB.Create(&models.Config{
			Title:      key,
			Content:    value,
			UpdateTime: time.Now(),
		})
		return rc.Error
	}
	return res.Error
}

func checkConfig(key string) bool {
	intCmd := redis.RedisClient.Exists(ctx, key)
	if intCmd.Val() == 1 {
		return true
	} else {
		return false
	}
}

func delConfig(key string) error {
	redis.RedisClient.Del(ctx, key)
	res := database.DB.Where(&models.Config{
		Title: key,
	}).Delete(models.Config{})
	return res.Error
}
