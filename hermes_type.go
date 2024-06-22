/*
Create: 2022/8/14
Project: Heimdallr
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"github.com/JJApplication/fushin/utils/json"
)

// hermes的操作

const (
	Hermes        = "Hermes"
	sendAlarmHtml = "sendAlarmHtml"
	sendAlarm     = "sendAlarm" // 纯文本类型不支持换行
)

type mailInfo struct {
	Type     string   `json:"type"`
	Message  string   `json:"message"`
	IsFile   bool     `json:"isFile"`
	Subject  string   `json:"subject"`
	Attach   []string `json:"attach"`
	To       []string `json:"to"`
	Cc       []string `json:"cc"`
	Bcc      []string `json:"bcc"`
	SyncTask bool     `json:"syncTask"`
	CronJob  string   `json:"cronJob"`
}

func mailReq(info mailInfo) string {
	d, err := json.Json.Marshal(info)
	if err != nil {
		return ""
	}
	return string(d)
}
