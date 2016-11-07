package database

import (
	"github.com/JesusIslam/salestock/configuration"
	"github.com/JesusIslam/salestock/logger"
	mgo "gopkg.in/mgo.v2"
)

var sess *mgo.Session

func init() {
	var err error
	sess, err = mgo.Dial(configuration.Database().URL)
	if err != nil {
		logger.Fatal("Failed to connect to database: ", err)
	}
}

func New() *mgo.Session {
	return sess.Clone()
}
