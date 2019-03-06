package main

import (
	"github.com/dy-dayan/community-srv-proposal/dal/db"
	"github.com/dy-dayan/community-srv-proposal/handler"
	"github.com/dy-dayan/community-srv-proposal/idl/dayan/community/srv-proposal"
	"github.com/dy-gopkg/kit"
	"github.com/dy-dayan/community-srv-proposal/util/config"
	"github.com/sirupsen/logrus"
)

func main() {
	kit.Init()

	// 初始化配置
	uconfig.Init()

	// 初始化数据库
	db.Init()

	//TODO 初始化缓存
	//cache.CacheInit()

	err := dayan_community_srv_proposal.RegisterProposalHandler(kit.DefaultService.Server(), &handler.Handler{})
	if err != nil {
		logrus.Fatalf("RegisterProposalHandler error:%v", err)
	}

	kit.Run()
}