package service

type ResourceStats struct {
	CpuUsagePercent uint8   `json:"cpu_used_percent"`
	UsedMemoryMB    uint64  `json:"used_memory_mb"`
	NetStat         NetStat `json:"network"`
}

type NetStat struct {
	SentKB uint64 `json:"sent_kb"`
	RecvKB uint64 `json:"recv_kb"`
}
