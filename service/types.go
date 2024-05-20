package service

type ResourceStats struct {
	CpuTime      *CpuTime `json:"cpu_usage_stat"`
	UsedMemoryMB uint64   `json:"used_memory_mb"`
	NetStat      *NetStat `json:"network"`
}

type NetStat struct {
	SentKB uint64 `json:"sent_kb"`
	RecvKB uint64 `json:"recv_kb"`
}

type CpuTime struct {
	Service uint64 `json:"application"`
	System  uint64 `json:"system"`
}
