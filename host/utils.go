package host

func Stats() (*ResourceStats, error) {
	resourceStats := &ResourceStats{}
	a, e := cpuUsage()
	if e != nil {
		return nil, e
	}
	resourceStats.CpuUsagePercent = a
	b, e := memoryStats()
	if e != nil {
		return nil, e
	}
	resourceStats.MemStat = *b
	c, e := diskStats()
	if e != nil {
		return nil, e
	}
	resourceStats.DiskStats = c
	d, e := netStats()
	if e != nil {
		return nil, e
	}
	resourceStats.NetStat = *d
	return resourceStats, nil
}
