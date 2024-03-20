package host

import (
	"math"
	"strings"

	"github.com/shirou/gopsutil/disk"
)

func diskStats() ([]DiskStat, error) {
	var partitions []disk.PartitionStat
	partitions, err := disk.Partitions(false)
	if err != nil {
		return nil, err
	}

	var diskStats = make([]DiskStat, 0)
	for _, value := range partitions {
		if strings.HasPrefix(value.Device, "/dev/loop") {
			continue
		} else if strings.HasPrefix(value.Mountpoint, "/var/lib/docker") {
			continue
		} else if strings.HasPrefix(value.Mountpoint, "/var/snap/") {
			continue
		}
		usageVal, err := disk.Usage(value.Mountpoint)
		if err != nil {
			continue
		}
		diskStats = append(diskStats, DiskStat{
			Path:       usageVal.Path,
			MountPoint: value.Device,
			TotalGB:    float32(math.Round(float64(usageVal.Total)/1024/1024/1024*100) / 100),
			UsedGB:     float32(math.Round(float64(usageVal.Used)/1024/1024/1024*100) / 100),
		})
	}

	return diskStats, nil
}
