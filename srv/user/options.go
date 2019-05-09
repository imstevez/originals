package main

import (
	"originals/sql"
	tokenProto "originals/srv/token/proto"
	"originals/srv/user/handler"
	"originals/srv/user/model"
	proto "originals/srv/user/proto"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"
)

// Initialize mysql database
func initMysqlDB(o *micro.Options) {
	o.BeforeStart = append(o.BeforeStart, func() error {
		log.Log("Initializing mysql database")
		err := sql.InitMysqlDB()
		if err != nil {
			log.Log(err)
			return err
		}

		return nil
	})
	o.AfterStop = append(o.AfterStop, func() error {
		log.Log("Close mysql database")
		if err := sql.MysqlDB.Close(); err != nil {
			log.Log("Close database failed: " + err.Error())
			return err
		}
		log.Log("Database closed")
		return nil
	})
}

// Register handler
func registerHandler(o *micro.Options) {
	o.BeforeStart = append(o.BeforeStart, func() error {
		log.Log("Register handler")
		tokenCli := tokenProto.NewTokenService("go.micro.srv.token", o.Client)
		if err := proto.RegisterUserHandler(o.Server,
			&handler.User{
				Model: &model.UserModel{
					DB: sql.MysqlDB,
				},
				TokenCli: tokenCli,
			}); err != nil {
			return err
		}
		return nil
	})
}
