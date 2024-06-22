/*
Create: 2022/8/20
Project: Heimdallr
Github: https://github.com/landers1037
Copyright Renj
*/

// Package main
package main

type SystemAlarmInfo struct {
	System       string
	Family       string
	Version      string
	Kernel       string
	BootTime     string
	CpuCount     int
	CpuPercent   float64
	MemUsed      float64
	MemAvail     float64
	ProcessCount int
}

type SystemAlarmAlert struct {
	MemUsed  float64
	MemAvail float64
}
