package server

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	minTCPPort = 0
	maxTCPPort = 65535
)

// PortRange innehåller en port range.
type PortRange struct {
	Min int16
	Max int16
}

// PortCheck innehåller information om man kunde pinga en port eller inte.
type PortCheck struct {
	Port    int16
	Success bool
}

// PortRangeCheck innehåller information om hela checken funkade eller inte.
type PortRangeCheck struct {
	Error string      `json:"error"`
	Ports []PortCheck `json:"ports"`
}

// Ping pingar en port range, finns flera olika syntax
// "3333" "22,3333" "22-80,3333"
func Ping(ip string, port string) PortRangeCheck {
	ports := []PortRange{}
	check := PortRangeCheck{Ports: []PortCheck{}}

	p := strings.Split(strings.Replace(port, " ", "", -1), ",")
	for _, port := range p {
		pSplit := strings.Split(port, "-")
		minPort, err := strconv.ParseUint(pSplit[0], 10, 16)
		if err != nil {
			return PortRangeCheck{Error: fmt.Sprintf("The value: %s can't be converted to an integer.", pSplit[0])}
		}
		if minPort < minTCPPort || minPort > maxTCPPort {
			return PortRangeCheck{Error: fmt.Sprintf("The value: %s is smaller or larger then the maximum TCP port range.", pSplit[0])}
		}
		maxPort := minPort

		if len(pSplit) == 2 {
			m, err := strconv.ParseUint(pSplit[1], 10, 16)
			if err != nil {
				return PortRangeCheck{Error: fmt.Sprintf("The value: %s can't be converted to an integer.", pSplit[1])}
			}
			if m < minTCPPort || m > maxTCPPort {
				return PortRangeCheck{Error: fmt.Sprintf("The value: %s is smaller or larger then the maximum TCP port range.", pSplit[1])}
			}
			maxPort = m
		}

		if minPort > maxPort {
			return PortRangeCheck{Error: fmt.Sprintf("Min value \"%d\" is larger then the max value \"%d\"", minPort, maxPort)}
		}
		ports = append(ports, PortRange{Min: (int16)(minPort), Max: (int16)(maxPort)})
	}
	for _, port := range ports {
		for i := port.Min; i <= port.Max; i++ {
			c := pingIP(ip, i)
			check.Ports = append(check.Ports, PortCheck{Port: i, Success: c})
			time.Sleep(time.Duration(1) * time.Second)
		}
	}
	return check
}

// pingIP pingar en specifik port på en ip
func pingIP(ip string, port int16) bool {
	conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), time.Duration(1)*time.Second)
	if err != nil {
		return false
	} else {
		conn.Close()
		return true
	}
}
