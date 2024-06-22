/*
Create: 2022/8/21
Project: Heimdallr
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"strings"

	"github.com/JJApplication/fushin/db/mongo"
)

// 告警信息推送

// 告警存储
// 数据库heimdallr会存储当前系统的告警信息
// 默认存储一个月

type Alarm struct {
	mongo.MetaModel `bson:",inline"`
	Title           string `json:"title" bson:"title"`
	Level           string `json:"level" bson:"level"`
	Message         string `json:"message" bson:"message"`
}

type AlarmBase struct {
	Title   string `json:"title" bson:"title"`
	Level   string `json:"level" bson:"level"`
	Message string `json:"message" bson:"message"`
}

func (alarm *Alarm) CollectionName() string {
	return "alarm"
}

const (
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
)

// CreateOneAlarm 信息以纯文本的方式存储
// mongo存储失败时不做任何事
func CreateOneAlarm(title, level, message string) error {
	logger.InfoF("start to push message to mongo [title]: %s, [level]: %s", title, level)
	return mongoC.Coll(&Alarm{}).Create(&Alarm{
		Title:   title,
		Level:   getLevel(level),
		Message: message,
	})
}

// PullAlarm 返回全部Alarm
func PullAlarm() []Alarm {
	var res []Alarm
	err := mongoC.Coll(&Alarm{}).SimpleFind(&res, mongo.M{})
	if err != nil {
		logger.ErrorF("pull all alarms from mongo error: %s", err.Error())
	}
	return res
}

// PurgeAlarm 删除Alarm
// 异步的删除 不返回任何错误
func PurgeAlarm() {

}

// PurgeAOneAlarm 删除Alarm
// 同步的删除 不返回任何错误
func PurgeAOneAlarm(id string) {

}

func getLevel(l string) string {
	switch strings.ToLower(l) {
	case LevelInfo:
		return LevelInfo
	case LevelWarn:
		return LevelWarn
	case LevelError:
		return LevelError
	default:
		return LevelInfo
	}
}
