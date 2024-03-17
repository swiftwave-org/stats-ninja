package service

import (
	"github.com/docker/docker/api/types"
	"math"
)

func calculateCPUPercentUnix(previousCPU, previousSystem uint64, v *types.StatsJSON) uint8 {
	var (
		cpuPercent  = 0.0
		cpuDelta    = float64(v.CPUStats.CPUUsage.TotalUsage) - float64(previousCPU)
		systemDelta = float64(v.CPUStats.SystemUsage) - float64(previousSystem)
	)
	if systemDelta > 0.0 && cpuDelta > 0.0 {
		cpuPercent = (cpuDelta / systemDelta) * 100.0
	}

	return uint8(math.Round(cpuPercent))
}

func memoryUsageMB(v *types.StatsJSON) uint64 {
	return uint64(math.Round(float64(v.MemoryStats.Usage) / 1024 / 1024))
}

func networkRecvKB(v *types.StatsJSON) uint64 {
	x := v.Networks
	var totalRecv float64
	for _, v := range x {
		totalRecv += float64(v.RxBytes)
	}
	return uint64(int(totalRecv/1024*100) / 100)
}

func networkSentKB(v *types.StatsJSON) uint64 {
	x := v.Networks
	var totalSent float64
	for _, v := range x {
		totalSent += float64(v.TxBytes)
	}
	return uint64(int(totalSent/1024*100) / 100)
}
