package service

import (
	"math"

	"github.com/docker/docker/api/types"
)

func calculateCPUPercentUnix(v *types.StatsJSON) uint8 {
	/*
		cpu_delta = cpu_stats.cpu_usage.total_usage - precpu_stats.cpu_usage.total_usage
		system_cpu_delta = cpu_stats.system_cpu_usage - precpu_stats.system_cpu_usage
		CPU usage % = (cpu_delta / system_cpu_delta) * 100.0
	*/
	cpuDelta := float64(v.CPUStats.CPUUsage.TotalUsage) - float64(v.PreCPUStats.CPUUsage.TotalUsage)
	systemDelta := float64(v.CPUStats.SystemUsage) - float64(v.PreCPUStats.SystemUsage)
	return uint8(math.Round((cpuDelta / systemDelta) * 100.0))
}

func memoryUsageMB(v *types.StatsJSON) uint64 {
	// used_memory = memory_stats.usage - memory_stats.stats.cache
	cache := uint64(0)
	if cacheStat, ok := v.MemoryStats.Stats["cache"]; ok {
		cache = cacheStat
	}
	return uint64(math.Round(float64(v.MemoryStats.Usage-cache) / 1024 / 1024))
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
