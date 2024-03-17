package service

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"io"
)

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
				CpuUsagePercent: 0,
				UsedMemoryMB:    0,
				NetStat: NetStat{
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
		rs.CpuUsagePercent = rs.CpuUsagePercent + calculateCPUPercentUnix(statsJSON.PreCPUStats.CPUUsage.TotalUsage, statsJSON.PreCPUStats.SystemUsage, &statsJSON)
		rs.UsedMemoryMB = rs.UsedMemoryMB + memoryUsageMB(&statsJSON)
		rs.NetStat.SentKB = rs.NetStat.SentKB + networkSentKB(&statsJSON)
	}

	return &statsMap, nil
}
