package main

import (
	"originals/sql"
	"originals/srv/user/handler"
	"originals/srv/user/model"
	"originals/srv/user/proto"

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
		if err := proto.RegisterUserSrvHandler(o.Server,
			&handler.UserSrvHandler{
				Model: &model.UserSrvModel{
					DB: sql.MysqlDB,
				},
			}); err != nil {
			return err
		}
		return nil
	})
}

// Say bye
func sayBye(o *micro.Options) {
	o.AfterStop = append(o.AfterStop, func() error {
		log.Log("Bye")
		return nil
	})
}
