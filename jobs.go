/*
Create: 2022/8/14
Project: Heimdallr
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"github.com/JJApplication/fushin/cron"
	"github.com/JJApplication/fushin/server/uds"
)

const (
	DISABLED = "disable"
)

// 后台定时执行的任务

func healthCheck() {
	if cf.JobHealthCheck == DISABLED {
		return
	}
	c := cron.NewGroup(cron.EveryFmt(cf.JobHealthCheck))
	_, err := c.AddFunc(func() {
		logger.Info("job [healthCheck] run")
	})
	if err != nil {
		logger.ErrorF("add job healthCheck error: %s", err.Error())
		return
	}
	c.Start()
}

// 每天早9点
func checkApps() {
	if cf.JobAppCheck == DISABLED {
		return
	}
	c := cron.NewGroup(cf.JobAppCheck)
	_, err := c.AddFunc(func() {
		logger.Info("job [checkApps] run")
		var allApps []string
		allApps = append(allApps, cf.Apps...)
		allApps = append(allApps, cf.ExtraApps...)
		appInfos, ok := checkProcess(allApps, true)
		if !ok {
			logger.Warn("job [checkApps] some apps not good")
		}

		res, err := udsc.SendWithRes(uds.Req{
			Operation: sendAlarmHtml,
			Data: mailReq(mailInfo{
				Type:     "",
				Message:  appAlarmInfo(appInfos),
				IsFile:   false,
				Subject:  TitleAppInfo,
				Attach:   nil,
				To:       []string{cf.To},
				Cc:       nil,
				Bcc:      nil,
				SyncTask: false,
				CronJob:  "",
			}),
			From: Heimdallr,
			To:   []string{Hermes},
		})
		if err != nil {
			logger.ErrorF("run job checkApps error: %s", err.Error())
		}
		if res.Error != "" {
			logger.ErrorF("res job checkApps error: %s", res.Error)
		}
	})
	if err != nil {
		logger.ErrorF("add job checkApps error: %s", err.Error())
		return
	}
	c.Start()
}

// 每天早上8点
func systemCheck() {
	if cf.JobSystemCheck == DISABLED {
		return
	}
	c := cron.NewGroup(cf.JobSystemCheck)
	_, err := c.AddFunc(func() {
		logger.Info("job [systemCheck] run")
		res, err := udsc.SendWithRes(uds.Req{
			Operation: sendAlarmHtml,
			Data: mailReq(mailInfo{
				Type:     "",
				Message:  systemAlarmInfo(),
				IsFile:   false,
				Subject:  TitleSystemInfo,
				Attach:   nil,
				To:       []string{cf.To},
				Cc:       nil,
				Bcc:      nil,
				SyncTask: false,
				CronJob:  "",
			}),
			From: Heimdallr,
			To:   []string{Hermes},
		})
		if err != nil {
			logger.ErrorF("run job systemCheck error: %s", err.Error())
		}
		if res.Error != "" {
			logger.ErrorF("res job systemCheck error: %s", res.Error)
		}
	})
	if err != nil {
		logger.ErrorF("add job systemCheck error: %s", err.Error())
		return
	}
	c.Start()
}

// 系统定时检查
// 暂时不使用 因为heimdallr运行时 cpu一定是高占用的
// 每6小时运行一次
// 内存基线60%
func systemLoopCheck() {
	if cf.JobSysLoopCheck == DISABLED {
		return
	}
	c := cron.NewGroup(cf.JobSysLoopCheck)
	_, err := c.AddFunc(func() {
		logger.Info("job [systemLoopCheck] run")
		memInfo := getMemUsed()
		if memInfo <= 0.5*100 {
			return
		}
		res, err := udsc.SendWithRes(uds.Req{
			Operation: sendAlarmHtml,
			Data: mailReq(mailInfo{
				Type:     "",
				Message:  systemAlarmAlert(),
				IsFile:   false,
				Subject:  TitleSysAlarm,
				Attach:   nil,
				To:       []string{cf.To},
				Cc:       nil,
				Bcc:      nil,
				SyncTask: false,
				CronJob:  "",
			}),
			From: Heimdallr,
			To:   []string{Hermes},
		})
		if err != nil {
			logger.ErrorF("run job systemLoopCheck error: %s", err.Error())
		}
		if res.Error != "" {
			logger.ErrorF("res job systemLoopCheck error: %s", res.Error)
		}
	})
	if err != nil {
		logger.ErrorF("add job systemLoopCheck error: %s", err.Error())
		return
	}
	c.Start()
}

// 服务检查 一小时一次
func checkAppsLoop() {
	if cf.JobAppLoopCheck == DISABLED {
		return
	}
	c := cron.NewGroup(cf.JobAppLoopCheck)
	_, err := c.AddFunc(func() {
		logger.Info("job [checkAppsLoop] run")
		var allApps []string
		allApps = append(allApps, cf.Apps...)
		allApps = append(allApps, cf.ExtraApps...)
		appInfos, ok := checkProcess(allApps, false)
		if !ok {
			logger.Warn("job [checkAppsLoop] some apps not good")
		}

		// 过滤app
		var badApps []appInfo
		for _, app := range appInfos {
			if app.Status == StatusBad {
				badApps = append(badApps, app)
			}
		}
		if len(badApps) <= 0 {
			return
		}
		res, err := udsc.SendWithRes(uds.Req{
			Operation: sendAlarmHtml,
			Data: mailReq(mailInfo{
				Type:     "",
				Message:  appAlarmAlert(badApps),
				IsFile:   false,
				Subject:  TitleAppAlarm,
				Attach:   nil,
				To:       []string{cf.To},
				Cc:       nil,
				Bcc:      nil,
				SyncTask: false,
				CronJob:  "",
			}),
			From: Heimdallr,
			To:   []string{Hermes},
		})
		if err != nil {
			logger.ErrorF("run job checkAppsLoop error: %s", err.Error())
		}
		if res.Error != "" {
			logger.ErrorF("res job checkAppsLoop error: %s", res.Error)
		}
	})
	if err != nil {
		logger.ErrorF("add job checkAppsLoop error: %s", err.Error())
		return
	}
	c.Start()
}
