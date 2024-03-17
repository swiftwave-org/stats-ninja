package host

type ResourceStats struct {
	CpuUsagePercent uint8      `json:"cpu_used_percent"`
	MemStat         MemoryStat `json:"memory"`
	DiskStats       []DiskStat `json:"disks"`
	NetStat         NetStat    `json:"network"`
}

type DiskStat struct {
	Path       string  `json:"path"`
	MountPoint string  `json:"mount_point"`
	TotalGB    float32 `json:"total_gb"`
	UsedGB     float32 `json:"used_gb"`
}

type MemoryStat struct {
	TotalGB  float32 `json:"total_gb"`
	UsedGB   float32 `json:"used_gb"`
	CachedGB float32 `json:"cached_gb"`
}

type NetStat struct {
	SentKB uint64 `json:"sent_kb"`
	RecvKB uint64 `json:"recv_kb"`
}
