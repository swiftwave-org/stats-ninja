package host

import (
	"github.com/shirou/gopsutil/net"
)

var lastCurrentKBytesSent uint64 = 0
var lastCurrentKBytesRecv uint64 = 0

func netStats() (*NetStat, error) {
	isPastStatFound := true
	if lastCurrentKBytesSent == 0 && lastCurrentKBytesRecv == 0 {
		isPastStatFound = false
	}
	netStats, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}
	if len(netStats) == 0 {
		return &NetStat{
			SentKB: 0,
			RecvKB: 0,
		}, nil
	}
	currentKBytesSent := uint64(netStats[0].BytesSent / 1024)
	currentKBytesRecv := uint64(netStats[1].BytesRecv / 1024)
	kiloBytesSentDiff := uint64(0)
	kiloBytesRecvDiff := uint64(0)
	if currentKBytesSent > lastCurrentKBytesSent {
		kiloBytesSentDiff = currentKBytesSent - lastCurrentKBytesSent
	}
	if currentKBytesRecv > lastCurrentKBytesRecv {
		kiloBytesRecvDiff = currentKBytesRecv - lastCurrentKBytesRecv
	}
	lastCurrentKBytesSent = currentKBytesSent
	lastCurrentKBytesRecv = currentKBytesRecv
	if !isPastStatFound {
		return &NetStat{
			SentKB: uint64(0),
			RecvKB: uint64(0),
		}, nil
	}
	return &NetStat{
		SentKB: kiloBytesSentDiff,
		RecvKB: kiloBytesRecvDiff,
	}, nil
}
