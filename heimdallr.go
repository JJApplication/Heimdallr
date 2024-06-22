/*
Create: 2022/8/6
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
	"github.com/JJApplication/fushin/utils/json"
)

var udsc *client.UDSClient
var logger *log.Logger
var mongoC *mongo.Mongo

func main() {
	logger = NewLogger()
	LoadConfig()
	udsc = NewClient()
	heimdallrServer := NewServer()
	heimdallrServer.AddFunc("ping", func(c *uds.UDSContext, req uds.Req) {
		_ = c.Response(uds.Res{
			Error: "",
			Data:  "",
			From:  Heimdallr,
			To:    nil,
		})
	})

	heimdallrServer.AddFunc("push", func(c *uds.UDSContext, req uds.Req) {
		var data AlarmBase
		if err := json.Json.UnmarshalFromString(req.Data, &data); err != nil {
			_ = c.Response(uds.Res{
				Error: err.Error(),
				Data:  "",
				From:  Heimdallr,
				To:    req.To,
			})
		} else {
			go func() {
				err := CreateOneAlarm(data.Title, data.Level, data.Message)
				if err != nil {
					logger.ErrorF("push message to mongo error: %s", err.Error())
				}
			}()
		}
	})

	heimdallrServer.AddFunc("pull", func(c *uds.UDSContext, req uds.Req) {
		res := PullAlarm()
		data, err := json.Json.MarshalToString(res)
		if err != nil {
			_ = c.Response(uds.Res{
				Error: err.Error(),
				Data:  "",
				From:  Heimdallr,
				To:    req.To,
			})
		} else {
			_ = c.Response(uds.Res{
				Error: "",
				Data:  data,
				From:  Heimdallr,
				To:    req.To,
			})
		}
	})

	heimdallrServer.AddFunc("purge", func(c *uds.UDSContext, req uds.Req) {
		var data Alarm
		if err := json.Json.UnmarshalFromString(req.Data, &data); err != nil {
			_ = c.Response(uds.Res{
				Error: err.Error(),
				Data:  "",
				From:  Heimdallr,
				To:    req.To,
			})
		} else {
			go func() {
				PurgeAOneAlarm(data.ID.String())
			}()
			_ = c.Response(uds.Res{
				Error: "",
				Data:  "",
				From:  Heimdallr,
				To:    req.To,
			})
		}
	})

	// 初始化数据库
	mongoC = NewMongo()
	err := mongoC.Init()
	if err != nil {
		logger.ErrorF("%s mongo client init error: %s", Heimdallr, err.Error())
	}
	// 初始化客户端
	err = udsc.Dial()
	if err != nil {
		logger.ErrorF("%s uds client dial error: %s", Heimdallr, err.Error())
	}
	// 初始化任务
	InitJobs()
	if err = heimdallrServer.Listen(); err != nil {
		logger.ErrorF("%s server start error: %s", Heimdallr, err.Error())
	}
}
