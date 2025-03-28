package sys_home

import (
	"gitee.com/mxcker/go-admin/server/core/response"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
	"github.com/shirou/gopsutil/v4/net"
	"runtime"
)

type controllerStruct struct{}

func (c *controllerStruct) Statistic(ctx *gin.Context) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	cpuInfo, err := cpu.Info()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	cpuPercents, err := cpu.Percent(0, true)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	// 获取内存信息
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	// 获取系统的平均负载
	loadAvg, err := load.Avg()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	// 获取所有分区的信息
	partitions, err := disk.Partitions(false)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	var disks []interface{}
	for _, part := range partitions {
		// 获取指定分区的使用情况
		if usageStat, err := disk.Usage(part.Mountpoint); err == nil && usageStat != nil {
			disks = append(disks, map[string]interface{}{
				"name":    part.Mountpoint,
				"total":   usageStat.Total / (1024 * 1024),
				"free":    usageStat.Free / (1024 * 1024),
				"used":    usageStat.Used / (1024 * 1024),
				"percent": usageStat.UsedPercent,
			})
		}
	}
	diskIO, err := disk.IOCounters()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	nets, err := net.Interfaces()
	if err != nil {
		response.Error(ctx, err)
		return
	}
	// 获取所有网络接口的信息
	netIO, err := net.IOCounters(true)
	if err != nil {
		response.Error(ctx, err)
		return
	}
	response.Success(ctx, gin.H{
		"go_threads": runtime.NumGoroutine(),
		"start_time": StartTime.UnixMilli(),
		"use_memory": memStats.Sys,
		"cpu_status": map[string]interface{}{
			"info":    cpuInfo,
			"counts":  len(cpuPercents),
			"percent": cpuPercents,
		},
		"mem_status": map[string]interface{}{
			"total":   float64(memInfo.Total) / (1024 * 1024),
			"free":    float64(memInfo.Free) / (1024 * 1024),
			"used":    float64(memInfo.Used) / (1024 * 1024),
			"percent": memInfo.UsedPercent,
		},
		"loadAvg": loadAvg,
		"disks":   disks,
		"diskIO":  diskIO,
		"nets":    nets,
		"netIO":   netIO,
	})
}
