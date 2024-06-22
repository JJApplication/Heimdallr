/*
Create: 2022/8/14
Project: Heimdallr
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"path/filepath"
	"strings"

	"github.com/JJApplication/fushin/utils/env"
	"github.com/JJApplication/octopus_meta"
)

// 配置 从环境变量加载

// v1 仅监控微服务的运行状态
// 默认监控jjapplication
// 基于octopus-meta 未发布的服务跳过
type heimdallrConfig struct {
	Apps        []string // 过滤掉未发布态后的jjapplication
	ExtraApps   []string // 额外监控的服务
	UnixAddress string   // 监听地址
	Talker      string   // 接收方地址
	To          string   // 收件人
	AppRoot     string   // app root
	MongoName   string
	MongoURL    string

	// 定时任务
	JobHealthCheck  string
	JobAppCheck     string
	JobSystemCheck  string
	JobSysLoopCheck string
	JobAppLoopCheck string
}

const (
	SymbolDot   = "."
	SymbolComma = ","
)

var cf heimdallrConfig

func init() {
	envLoader := env.EnvLoader{}
	cf = heimdallrConfig{
		Apps:        getApps(envLoader.Get("APP_ROOT").Raw()),
		ExtraApps:   strings.Split(strings.TrimSpace(envLoader.Get("ExtraApps").Raw()), SymbolComma),
		UnixAddress: envLoader.Get("UnixAddress").Raw(),
		Talker:      envLoader.Get("Talker").Raw(),
		To:          envLoader.Get("To").Raw(),
		AppRoot:     envLoader.Get("APP_ROOT").Raw(),
		MongoName:   envLoader.Get("MongoName").Raw(),
		MongoURL:    envLoader.Get("MongoURL").Raw(),

		// jobs
		JobHealthCheck:  envLoader.Get("").MustString("30m"),
		JobAppCheck:     envLoader.Get("").MustString("0 0 9 * * ?"),
		JobSystemCheck:  envLoader.Get("").MustString("0 0 8 * * ?"),
		JobSysLoopCheck: envLoader.Get("").MustString("0 0 0/6 * * ?"),
		JobAppLoopCheck: envLoader.Get("").MustString("0 0 0/1 * * ?"),
	}
}

func getApps(appRoot string) []string {
	var apps []string
	appsMap, err := octopus_meta.AutoLoad()
	if err != nil {
		return apps
	}
	// 默认跳过NoEngine NoEngine异常时服务才会异常
	for k, v := range appsMap {
		if v.ReleaseStatus == octopus_meta.Published &&
			v.Type != octopus_meta.TypeNoEngine &&
			v.Type != octopus_meta.TypeContainer {
			apps = append(apps, filepath.Join(appRoot, k))
		}
	}
	return apps
}
