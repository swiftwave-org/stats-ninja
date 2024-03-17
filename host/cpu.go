package host

import (
	"math"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

func cpuUsage() (uint8, error) {
	cpuRates, err := cpu.Percent(time.Second, true)
	if err != nil {
		return 0, err
	}
	// calculate average cpu usage
	var totalCpuUsage float64
	for _, rate := range cpuRates {
		totalCpuUsage += rate
	}
	// round to ceil
	cpuUsage := uint8(math.Round(totalCpuUsage / float64(len(cpuRates))))
	return cpuUsage, nil
}
