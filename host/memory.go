package host

import (
	"math"

	"github.com/shirou/gopsutil/mem"
)

func memoryStats() (*MemoryStat, error) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	return &MemoryStat{
		TotalGB:  float32(math.Round(float64(memory.Total)/1024/1024/1024*100) / 100),
		UsedGB:   float32(math.Round(float64(memory.Used)/1024/1024/1024*100) / 100),
		CachedGB: float32(math.Round(float64(memory.Cached)/1024/1024/1024*100) / 100),
	}, nil
}
