package host

import (
	"errors"

	"github.com/shirou/gopsutil/net"
)

var lastCurrentBytesSent uint64 = 0
var lastCurrentBytesRecv uint64 = 0

func netStats() (*NetStat, error) {
	isPastStatFound := true
	if lastCurrentBytesSent == 0 && lastCurrentBytesRecv == 0 {
		isPastStatFound = false
	}
	netStats, err := net.IOCounters(false)
	if err != nil {
		return nil, err
	}
	IO := make(map[string][]uint64)
	for _, IOStat := range netStats {
		nic := []uint64{IOStat.BytesSent, IOStat.BytesRecv}
		IO[IOStat.Name] = nic
	}
	if len(IO) == 0 {
		return &NetStat{
			SentKB: 0,
			RecvKB: 0,
		}, nil
	}
	if _, ok := IO["all"]; !ok {
		return nil, errors.New("interface not found")
	}
	allNet := IO["all"]
	currentBytesSent := allNet[0]
	currentBytesRecv := allNet[1]
	bytesSent := currentBytesSent - lastCurrentBytesSent
	bytesRecv := currentBytesRecv - lastCurrentBytesRecv
	lastCurrentBytesSent = currentBytesSent
	lastCurrentBytesRecv = currentBytesRecv
	if !isPastStatFound {
		return &NetStat{
			SentKB: 0,
			RecvKB: 0,
		}, nil
	}
	if bytesSent <= 0 {
		bytesSent = 0
	}
	if bytesRecv <= 0 {
		bytesRecv = 0
	}
	return &NetStat{
		SentKB: bytesSent / 1024,
		RecvKB: bytesRecv / 1024,
	}, nil
}
