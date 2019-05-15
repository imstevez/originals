package main

import (
	"originals/redis"
	"originals/srv/token/handler"
	"originals/srv/token/model"
	"originals/srv/token/proto"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
)

// 初始化redis
func initRedis(o *micro.Options) {
	o.BeforeStart = append(o.BeforeStart, func() error {
		log.Log("Initializing redis")
		err := redis.InitRedis()
		if err != nil {
			log.Log(err)
			return err
		}
		return nil
	})
	o.AfterStop = append(o.AfterStop, func() error {
		log.Log("Close redis")
		if err := redis.Redis.Close(); err != nil {
			log.Log("Close redis failed: " + err.Error())
			return err
		}
		log.Log("redis closed")
		return nil
	})
}

// 注册handler
func registerHandler(o *micro.Options) {
	o.BeforeStart = append(o.BeforeStart, func() error {
		log.Log("Register token service handler")
		if err := proto.RegisterTokenHandler(o.Server,
			&handler.Token{
				Model: &model.TokenModel{
					Redis: redis.Redis,
				},
			}); err != nil {
			return err
		}
		return nil
	})
}
