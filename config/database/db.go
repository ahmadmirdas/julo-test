package database

import (
	pg "github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
)

type ParamConn struct {
	Username    string
	Password    string
	Host        string
	Port        string
	Database    string
	MaxConn     int
	MinIdleConn int
	MaxRetries  int
}

func DbConn(param ParamConn) *pg.DB {
	logrus.Info("Initialize Database")

	db := pg.Connect(&pg.Options{
		User:         param.Username,
		Password:     param.Password,
		Addr:         param.Host + ":" + param.Port,
		Database:     param.Database,
		PoolSize:     param.MaxConn, // NOTE: increase this if you want to handle larger/bigger concurrent request
		MinIdleConns: param.MinIdleConn,
		MaxRetries:   param.MaxRetries,
	})

	return db
}
