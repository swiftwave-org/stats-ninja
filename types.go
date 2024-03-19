package main

import (
	"encoding/json"
	"github.com/swiftwave-org/stats_ninja/host"
	"github.com/swiftwave-org/stats_ninja/service"
)

type StatsData struct {
	SystemStat   host.ResourceStats                `json:"system"`
	ServiceStats map[string]*service.ResourceStats `json:"services"`
	TimeStamp    uint64                            `json:"timestamp"`
}

func (s *StatsData) JSON() ([]byte, error) {
	// convert to json
	jsonData, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
