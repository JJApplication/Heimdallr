/*
Create: 2022/9/3
Project: Heimdallr
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"context"

	"github.com/JJApplication/fushin/db/mongo"
	"github.com/JJApplication/octopus_meta"
)

type App struct {
	Meta octopus_meta.App `json:"meta" bson:"meta"`
}

type DaoAPP struct {
	mongo.MetaModel `bson:",inline"`
	App             `json:"app" bson:"app"`
}

func (app *DaoAPP) CollectionName() string {
	return "microservice"
}

// 检查是否是手动停止
func isStopByApollo(app string) bool {
	var data DaoAPP
	err := mongoC.Coll(&DaoAPP{}).FindOne(context.Background(), mongo.M{"app.meta.name": app}).Decode(&data)
	if err != nil {
		logger.ErrorF("seek app meta error: %s", err.Error())
		return false
	}
	return data.Meta.Runtime.StopOperation
}
