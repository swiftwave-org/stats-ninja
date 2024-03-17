package main

import (
	"github.com/docker/docker/client"
	"github.com/swiftwave-org/stats_ninja/host"
	"github.com/swiftwave-org/stats_ninja/service"
	"time"
)

func fetchStats(dockerClient *client.Client) (*StatsData, error) {
	// fetch system stats
	systemStats, err := host.Stats()
	if err != nil {
		return nil, err
	}
	// fetch service stats
	serviceStats, err := service.Stats(dockerClient)
	if err != nil {
		return nil, err
	}
	// record timestamp
	unixTimestamp := uint64(time.Now().Unix())
	// make it to nearest minute (floor)
	unixTimestamp = unixTimestamp - (unixTimestamp % 60)
	// create a new StatsData
	statsData := &StatsData{
		SystemStat:   *systemStats,
		ServiceStats: *serviceStats,
		TimeStamp:    unixTimestamp,
	}
	return statsData, nil
}
