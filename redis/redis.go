package redis

import "github.com/go-redis/redis"

var Redis *redis.Client

func InitRedis() error {
	option := &redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "root",
		DB:       0,
		PoolSize: 5,
	}
	Redis := redis.NewClient(option)
	if _, err := Redis.Ping().Result(); err != nil {
		return err
	}
	return nil
}
