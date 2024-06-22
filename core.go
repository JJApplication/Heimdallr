/*
Create: 2022/8/14
Project: Heimdallr
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	client "github.com/JJApplication/fushin/client/uds"
	"github.com/JJApplication/fushin/db/mongo"
	"github.com/JJApplication/fushin/log"
	"github.com/JJApplication/fushin/server/uds"
)

// NewServer 新建uds服务器 用于心跳
func NewServer() *uds.UDSServer {
	s := uds.Default(cf.UnixAddress)
	s.Option.AutoRecover = true
	s.Option.AutoCheck = false
	s.Option.MaxSize = 5 << 20
	logger.InfoF("%s uds server run @ [%s]", Heimdallr, s.Name)
	return s
}

// NewClient 新建uds客户端
func NewClient() *client.UDSClient {
	logger.InfoF("%s uds client dial @ [%s]", Heimdallr, cf.Talker)
	return &client.UDSClient{
		Addr:        cf.Talker,
		MaxRecvSize: 1 << 20,
	}
}

func NewLogger() *log.Logger {
	return log.Default(Heimdallr)
}

func LoadConfig() {
	logger.Info("config loaded from env")
	logger.InfoF("config: %+v", cf)
}

func NewMongo() *mongo.Mongo {
	m := &mongo.Mongo{
		ContextTimeout: 10,
		DBName:         cf.MongoName,
		URL:            cf.MongoURL,
	}
	return m
}

// InitJobs 初始化定时任务
func InitJobs() {
	healthCheck()
	checkApps()
	systemCheck()
	systemLoopCheck()
	checkAppsLoop()
}
