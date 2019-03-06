package uconfig

import (
	"github.com/sirupsen/logrus"
	"github.com/micro/go-config"
)

type (
	MongoDBConfig struct {
		Addr []string `json:"addr"`
		Username string `json:"username"`
		Password string `json:"password"`
		PoolLimit int `json:"poolLimit"`
	}
)

var (
	DefaultMgoConf MongoDBConfig
)

func Init() {
	// 加载mongo配置
	err := config.Get("mongodb").Scan(&DefaultMgoConf)
	if err != nil {
		logrus.Fatalf("get mgo config error: %s", err)
	}

	if len(DefaultMgoConf.Addr) == 0 {
		logrus.Fatalf("invalid mgo addr")
	}

}