package service

import (
	"context"
	"encoding/json"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// store previous stats
var serviceLastNetStats map[string]*NetStat

func init() {
	serviceLastNetStats = make(map[string]*NetStat)
}

func Stats(cl *client.Client) (*map[string]*ResourceStats, error) {
	// map of service name to stats
	statsMap := make(map[string]*ResourceStats)

	// fetch all the containers
	containers, err := cl.ContainerList(context.Background(), container.ListOptions{
		Size: false,
		All:  false,
		Filters: filters.NewArgs(
			filters.Arg("status", "running"),
		),
	})
	if err != nil {
		return nil, err
	}
	// iterate over the containers
	for _, c := range containers {
		// fetch stats
		stats, err := cl.ContainerStats(context.Background(), c.ID, false)
		if err != nil {
			continue
		}
		// ignore standalone containers
		if serviceName, ok := c.Labels["com.docker.swarm.service.name"]; !ok || serviceName == "" {
			continue
		}
		// get service name
		serviceName := c.Labels["com.docker.swarm.service.name"]
		// create a new ResourceStats if it doesn't exist
		if _, ok := statsMap[serviceName]; !ok {
			statsMap[serviceName] = &ResourceStats{
				ServiceCpuTime: 0,
				SystemCpuTime:  0,
				UsedMemoryMB:   0,
				NetStat: &NetStat{
					SentKB: 0,
					RecvKB: 0,
				},
			}
		}
		// fetch the ResourceStats
		rs := statsMap[serviceName]
		// Read the content of rc into a byte slice
		content, err := io.ReadAll(stats.Body)
		if err != nil {
			continue
		}
		// convert to types.StatsJSON
		var statsJSON types.StatsJSON
		err = json.Unmarshal(content, &statsJSON)
		if err != nil {
			continue
		}

		// save the stats
		rs.ServiceCpuTime = rs.ServiceCpuTime + uint64(statsJSON.CPUStats.CPUUsage.TotalUsage-statsJSON.PreCPUStats.CPUUsage.TotalUsage)
		rs.SystemCpuTime = uint64(statsJSON.CPUStats.SystemUsage - statsJSON.PreCPUStats.SystemUsage)
		rs.UsedMemoryMB = rs.UsedMemoryMB + memoryUsageMB(&statsJSON)
		rs.NetStat.SentKB = rs.NetStat.SentKB + networkSentKB(&statsJSON)
		rs.NetStat.RecvKB = rs.NetStat.RecvKB + networkRecvKB(&statsJSON)
	}

	calculateNetStatDiffFromLastRecord(&statsMap)

	return &statsMap, nil
}

func calculateNetStatDiffFromLastRecord(statsMap *map[string]*ResourceStats) {
	statsMapRef := *statsMap

	// iterate over the statsMap
	for serviceName, stats := range statsMapRef {
		// check if old stats exist
		if _, ok := serviceLastNetStats[serviceName]; !ok {
			// add this as the last stats
			serviceLastNetStats[serviceName] = &NetStat{
				SentKB: stats.NetStat.SentKB,
				RecvKB: stats.NetStat.RecvKB,
			}
			// mark the current net stats as 0
			statsMapRef[serviceName].NetStat.SentKB = 0
			statsMapRef[serviceName].NetStat.RecvKB = 0
		} else {
			oldSentKB := serviceLastNetStats[serviceName].SentKB
			oldRecvKB := serviceLastNetStats[serviceName].RecvKB
			currentSentKB := stats.NetStat.SentKB
			currentRecvKB := stats.NetStat.RecvKB
			if currentSentKB < oldSentKB {
				statsMapRef[serviceName].NetStat.SentKB = 0
			} else {
				statsMapRef[serviceName].NetStat.SentKB = currentSentKB - oldSentKB
			}
			if currentRecvKB < oldRecvKB {
				statsMapRef[serviceName].NetStat.RecvKB = 0
			} else {
				statsMapRef[serviceName].NetStat.RecvKB = currentRecvKB - oldRecvKB
			}
			serviceLastNetStats[serviceName].SentKB = currentSentKB
			serviceLastNetStats[serviceName].RecvKB = currentRecvKB
		}
	}
}
