/*
Create: 2022/8/14
Project: Heimdallr
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

import (
	"math"
	"path/filepath"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/load"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/process"
)

// 调用操作系统api

const (
	StatusBad  = "bad"
	StatusGood = "good"
)

// 获取cpu核心数
func getCpuCount() int {
	c, err := cpu.Counts(true)
	if err != nil {
		return 0
	}
	return c
}

// 获取cpu占用率
// 默认只计算单核
func getCpu() float64 {
	data, err := cpu.Percent(time.Second, true)
	if err != nil {
		logger.ErrorF("get cpu percent error: %s", err.Error())
		return 0
	}
	// 计算总值
	var all float64
	for _, v := range data {
		all += v
	}
	return math.Floor(all / float64(len(data)))
}

// 获取系统负载
func getLoad() string {
	info, err := load.Avg()
	if err != nil {
		return ""
	}
	return info.String()
}

// 获取内存使用率
// 即used/total
func getMemUsed() float64 {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return v.UsedPercent
}

// 获取内存空闲
func getMemAvail() float64 {
	v, err := mem.VirtualMemory()
	if err != nil {
		return 0
	}
	return float64(v.Available)
}

// 获取磁盘信息

// 获取主机启动信息
func getBoot() string {
	start, err := host.BootTime()
	if err != nil {
		return ""
	}
	t := time.Unix(int64(start), 0)
	return t.Local().Format("2006-01-02 15:04:05")
}

// 获取内核信息
func getKernel() string {
	k, err := host.KernelVersion()
	if err != nil {
		return ""
	}
	return k
}

// 获取平台信息
// platform family version
func getPlatform() (string, string, string) {
	p, f, v, err := host.PlatformInformation()
	if err != nil {
		return "", "", ""
	}
	return p, f, v
}

// 获取进程数
func getProcessCount() int {
	ps, err := process.Processes()
	if err != nil {
		return 0
	}
	return len(ps)
}

// 获取当前系统进程
func getProcess() []*process.Process {
	ps, err := process.Processes()
	if err != nil {
		return ps
	}
	return ps
}

// 获取进程信息
// p.String 为pid信息
func getProcessInfo(p *process.Process) string {
	return p.String()
}

type appInfo struct {
	App           string
	Cmd           string // 可执行程序命令
	Status        string // Linux的程序运行态
	Pid           int32
	CreateTime    int64
	CpuPercent    float64
	MemPercent    float32
	MemRss        uint64
	Connections   int
	NumberThreads int32
	IOReadCount   uint64
	IOWriteCount  uint64
	IOReadBytes   uint64
	IOWriteBytes  uint64
}

// 检查给定app的进程信息
// 有错误时 为false
func checkProcess(apps []string, all bool) ([]appInfo, bool) {
	ps := getProcess()
	if len(ps) == 0 {
		return nil, false
	}
	var infos []appInfo
	var down = true
	// 对于子进程 暂不处理
	// 跳过空的app
	for _, app := range apps {
		if app == "" {
			continue
		}
		if !all && isStopByApollo(filepath.Base(app)) {
			continue
		}
		tmpInfo := appInfo{App: filepath.Base(app), Status: StatusBad}
		for _, p := range ps {
			pName, _ := p.Cmdline()
			if strings.Contains(pName, app) {
				// cpu
				cpuP, _ := p.CPUPercent()
				memP, _ := p.MemoryPercent()
				memInfo, _ := p.MemoryInfo()
				createTime, _ := p.CreateTime()
				status, _ := p.Status()
				conns, _ := p.Connections()
				numThreads, _ := p.NumThreads()
				io, _ := p.IOCounters()
				tmpInfo = appInfo{
					App:           filepath.Base(app),
					Cmd:           pName,
					Status:        status,
					Pid:           p.Pid,
					CreateTime:    createTime,
					CpuPercent:    cpuP,
					MemPercent:    memP,
					MemRss:        memInfo.RSS,
					Connections:   len(conns),
					NumberThreads: numThreads,
					IOReadCount:   io.ReadCount,
					IOWriteCount:  io.WriteCount,
					IOReadBytes:   io.ReadBytes,
					IOWriteBytes:  io.WriteBytes,
				}
				break
			}
			down = false
		}
		infos = append(infos, tmpInfo)
	}
	return infos, down
}
