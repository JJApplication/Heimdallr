/*
Create: 2022/8/17
Project: Heimdallr
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

type Message struct {
	ID         string `json:"id" bson:"id"`                   // uuid
	CreateTime int64  `json:"create_time" bson:"create_time"` // 时间戳 精确到s
	Level      string `json:"level" bson:"level"`             // 告警级别
	Info       string `json:"info" bson:"info"`
	Source     string `json:"source" bson:"source"` // 告警来源
}

// 告警信息的来源
// 按类型分为服务类 系统类 容器类 其他
const (
	SourceApp       = "app"
	SourceSystem    = "system"
	SourceContainer = "container"
	SourceOther     = "other"
)
