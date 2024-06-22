/*
Create: 2022/8/20
Project: Heimdallr
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

// 告警信息模板

const (
	SystemInfoTemplate = "<h4>环境信息</h4>操作系统: %s</br>系统家族: %s</br>系统版本: %s</br>内核版本: %s</br>上次启动时间: %s</br>CPU核心数: %d</br>CPU使用率: %v</br>内存使用率: %v</br>可用内存: %v</br>运行进程数: %d</br>"

	SystemAlertTemplate = "<h4>系统告警</h4><p>内存占用达到阈值</p></br><p>内存占用: %f%%</br>空闲内存: %f bytes</p>"
)

const (
	SystemInfoTemplateMongo  = "环境信息\n操作系统: %s\n系统家族: %s\n系统版本: %s\n内核版本: %s\n上次启动时间: %s\nCPU核心数: %d\nCPU使用率: %v\n内存使用率: %v\n可用内存: %v\n运行进程数: %d"
	SystemAlertTemplateMongo = "系统告警\n内存占用达到阈值\n内存占用: %.3f%%\n空闲内存: %s"
)
