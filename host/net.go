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
	IO := make(map[string][]int32)
	for _, IOStat := range netStats {
		nic := []int32{int32(IOStat.BytesSent), int32(IOStat.BytesRecv)}
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
	currentBytesSent := uint64(allNet[0])
	currentBytesRecv := uint64(allNet[1])
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
		SentKB: uint64(bytesSent / 1024),
		RecvKB: uint64(bytesRecv / 1024),
	}, nil
}
