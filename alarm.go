/*
Create: 2022/8/14
Project: Heimdallr
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"fmt"
	"strconv"
	"time"
)

// 生成告警信息

const (
	TitleSystemInfo = "[SYSTEM INFO]系统状态"
	TitleAppInfo    = "[APP INFO]服务状态"
	TitleAppAlarm   = "[APP ALARM]服务告警"
	TitleSysAlarm   = "[SYSTEM ALARM]系统告警"
)

// 生成系统信息
// 非告警 定时发送
func systemAlarmInfo() string {
	p, f, v := getPlatform()
	sysInfo := SystemAlarmInfo{
		System:       p,
		Family:       f,
		Version:      v,
		Kernel:       getKernel(),
		BootTime:     getBoot(),
		CpuCount:     getCpuCount(),
		CpuPercent:   getCpu(),
		MemUsed:      getMemUsed(),
		MemAvail:     getMemAvail(),
		ProcessCount: getProcessCount(),
	}
	go pushSystemInfo(sysInfo)
	return fmt.Sprintf(SystemInfoTemplate,
		sysInfo.System,
		sysInfo.Family,
		sysInfo.Version,
		sysInfo.Kernel,
		sysInfo.BootTime,
		sysInfo.CpuCount,
		fmt.Sprintf("%.3f%%", sysInfo.CpuPercent),
		fmt.Sprintf("%.3f%%", sysInfo.MemUsed),
		fmt.Sprintf("%s", calcSize(int64(sysInfo.MemAvail))),
		sysInfo.ProcessCount,
	)
}

func systemAlarmAlert() string {
	sysAlert := SystemAlarmAlert{
		MemUsed:  getMemUsed(),
		MemAvail: getMemAvail(),
	}
	go pushSystemAlert(sysAlert)
	return fmt.Sprintf(SystemAlertTemplate, sysAlert.MemUsed, sysAlert.MemAvail)
}

// 服务监控定时信息
func appAlarmInfo(appInfos []appInfo) string {
	header := `<h4 style="color: #378de5">微服务状态</h4>`
	body := ""
	for _, app := range appInfos {
		if app.Status == StatusBad {
			appBody := `<div style="font-size: 0.85rem;margin-bottom: 0.5rem">%s</div>`
			pre := fmt.Sprintf(`<strong style="font-size: 1rem;color: #dc4905">[*] %s</strong>`, app.App)
			infoBody := fmt.Sprintf(appBody, pre)
			body = body + infoBody
			continue
		}
		appBody := `<div style="font-size: 0.85rem;margin-bottom: 0.5rem">%s%s%s%s%s%s%s%s%s</div>`
		pre := fmt.Sprintf(`<strong style="font-size: 1rem;color: #30a24c">[*] %s</strong>`, app.App)
		t := fmt.Sprintf(`<p style="margin: 4px 0"><strong>创建时间:</strong> %s</p>`, time.Unix(app.CreateTime/1000, 0).Local().Format("2006-01-02 15:04:05"))
		pid := fmt.Sprintf(`<p style="margin: 4px 0"><strong>PID:</strong> %d</p>`, app.Pid)
		cpu := fmt.Sprintf(`<p style="margin: 4px 0"><strong>CPU使用:</strong> %.3f%%</p>`, app.CpuPercent)
		mem := fmt.Sprintf(`<p style="margin: 4px 0"><strong>内存使用:</strong> %.3f%%</p>`, app.MemPercent)
		memRss := fmt.Sprintf(`<p style="margin: 4px 0"><strong>内存占用:</strong> %s</p>`, calcSize(int64(app.MemRss)))
		conn := fmt.Sprintf(`<p style="margin: 4px 0"><strong>连接数:</strong> %d</p>`, app.Connections)
		thread := fmt.Sprintf(`<p style="margin: 4px 0"><strong>线程数:</strong> %d</p>`, app.NumberThreads)
		io := fmt.Sprintf(`<p style="margin: 4px 0"><strong>IO使用:</strong> <span style="padding: 0 4px">读次数: %d</span><span style="padding: 0 4px">写次数: %d</span><span style="padding: 0 4px">读字节: %d</span><span style="padding: 0 4px">写字节: %d</span></p>`,
			app.IOReadCount,
			app.IOWriteCount,
			app.IOReadBytes,
			app.IOWriteBytes,
		)

		infoBody := fmt.Sprintf(appBody, pre, t, pid, cpu, mem, memRss, conn, thread, io)
		body = body + infoBody
	}
	go pushAppInfo(appInfos)
	return fmt.Sprintf("%s%s", header, body)
}

// 服务监控告警信息
// 上报出错的微服务
func appAlarmAlert(appInfos []appInfo) string {
	header := `<h4 style="color: #dc4905">微服务状态异常</h4>`
	body := ""
	for _, app := range appInfos {
		if app.Status == StatusBad {
			appBody := `<div style="font-size: 0.85rem;margin-bottom: 0.5rem">%s</div>`
			pre := fmt.Sprintf(`<strong style="font-size: 1rem;color: #dc4905">[*] %s</strong>`, app.App)
			infoBody := fmt.Sprintf(appBody, pre)
			body = body + infoBody
			continue
		}
	}
	go pushAppAlert(appInfos)
	return fmt.Sprintf("%s%s", header, body)
}

const (
	KB = 1024
	MB = 2 << 20
	GB = 2 << 30
)

// calcSize 计算工具类
func calcSize(s int64) string {
	if s == 0 {
		return "0kb"
	} else if s < KB {
		return fmt.Sprintf("%sb", strconv.FormatInt(s, 10))
	} else if s >= KB && s < MB {
		return fmt.Sprintf("%skb", strconv.FormatInt(s/KB, 10))
	} else if s >= MB && s < GB {
		return fmt.Sprintf("%smb", strconv.FormatInt(s/MB, 10))
	} else {
		return fmt.Sprintf("%sgb", strconv.FormatInt(s/GB, 10))
	}
}
