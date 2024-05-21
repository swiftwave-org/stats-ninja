package service

type ResourceStats struct {
	ServiceCpuTime uint64   `json:"service_cpu_time"`
	SystemCpuTime  uint64   `json:"system_cpu_time"`
	UsedMemoryMB   uint64   `json:"used_memory_mb"`
	NetStat        *NetStat `json:"network"`
}

type NetStat struct {
	SentKB uint64 `json:"sent_kb"`
	RecvKB uint64 `json:"recv_kb"`
}
