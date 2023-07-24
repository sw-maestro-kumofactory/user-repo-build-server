package config

import (
	"github.com/go-redis/redis"
)

func InitRedisClient() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "kumo-deploy-test.qecpup.ng.0001.apn2.cache.amazonaws.com:6379",
		Password: "",
		DB:       0,
	})

	_, err := client.Ping().Result()

	return client, err
}
