package redis

import "github.com/go-redis/redis"

var (
	Redis *redis.Client
	Nil   = redis.Nil
)

func InitRedis() error {
	option := &redis.Options{
		Addr:     "www.koogo.net:6379",
		Password: "root",
		DB:       0,
		PoolSize: 5,
	}
	Redis = redis.NewClient(option)
	if _, err := Redis.Ping().Result(); err != nil {
		return err
	}
	return nil
}
