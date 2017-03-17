package client

import (
	"github.com/shirou/gopsutil/net"
)

// NetResponse innehåller information on klienten's nätverks I/O
type NetResponse struct {
	Error string `json:"error"`
	Sent  uint64 `json:"sent"`
	Recv  uint64 `json:"recv"`
}

// UptimeCheck hämtar klienten's nätverks I/O
func NetworkCheck(cmd Command) NetResponse {
	i, err := net.IOCounters(false)
	if err != nil {
		return NetResponse{Error: err.Error()}
	}
	if len(i) <= 0 {
		return NetResponse{Error: "Can't get network I/O"}
	}

	io := i[0]

	return NetResponse{
		Sent: io.BytesSent,
		Recv: io.BytesRecv,
	}
}
