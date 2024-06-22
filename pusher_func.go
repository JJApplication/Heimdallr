/*
Create: 2022/8/21
Project: Heimdallr
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"fmt"
)

// 存储上报的通知

func pushSystemInfo(sysInfo SystemAlarmInfo) {
	body := fmt.Sprintf(SystemInfoTemplateMongo,
		sysInfo.System,
		sysInfo.Family,
		sysInfo.Version,
		sysInfo.Kernel,
		sysInfo.BootTime,
		sysInfo.CpuCount,
		fmt.Sprintf("%.3f%%", sysInfo.CpuPercent),
		fmt.Sprintf("%.3f%%", sysInfo.MemUsed),
		fmt.Sprintf("%s", calcSize(int64(sysInfo.MemAvail))),
		sysInfo.ProcessCount)
	err := CreateOneAlarm(TitleSystemInfo, LevelInfo, body)
	if err != nil {
		logger.ErrorF("[SystemInfo] push message to mongo error: %s", err.Error())
	}
}

func pushSystemAlert(sysAlert SystemAlarmAlert) {
	body := fmt.Sprintf(SystemAlertTemplateMongo, sysAlert.MemUsed, calcSize(int64(sysAlert.MemAvail)))
	err := CreateOneAlarm(TitleSysAlarm, LevelError, body)
	if err != nil {
		logger.ErrorF("[SystemAlert] push message to mongo error: %s", err.Error())
	}
}

func pushAppInfo(appInfos []appInfo) {
	var hasBad bool
	for _, app := range appInfos {
		if app.Status == StatusBad {
			hasBad = true
			break
		}
	}
	if hasBad {
		body := "微服务状态正常"
		err := CreateOneAlarm(TitleAppAlarm, LevelInfo, body)
		if err != nil {
			logger.ErrorF("[AppInfo] push message to mongo error: %s", err.Error())
		}
	}
}

func pushAppAlert(appInfos []appInfo) {
	body := "微服务状态异常\n"
	for _, app := range appInfos {
		if app.Status == StatusBad {
			appBody := fmt.Sprintf("[*] %s\n", app.App)
			body = body + appBody
			continue
		}
	}
	err := CreateOneAlarm(TitleAppAlarm, LevelError, body)
	if err != nil {
		logger.ErrorF("[AppAlert] push message to mongo error: %s", err.Error())
	}
}
