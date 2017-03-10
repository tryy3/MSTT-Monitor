package client

import "github.com/shirou/gopsutil/host"
import "github.com/shirou/gopsutil/net"

type UptimeResponse struct {
	Error  string `json:"error"`
	Uptime uint64 `json:"uptime"`
}

func UptimeCheck(cmd Command) UptimeResponse {
	up, err := host.Uptime()
	if err != nil {
		return UptimeResponse{Error: err.Error()}
	}
	return UptimeResponse{Uptime: up}
}

type InterfaceResponse struct {
	Name string   `json:"name"`
	IPs  []string `json:"ips"` // Kan ha både ipv4 och ipv6
}

type InfoResponse struct {
	Error      string              `json:"error"`
	Hostname   string              `json:"hostname"`
	OS         string              `json:"os"`
	Platform   string              `json:"platform"`
	Interfaces []InterfaceResponse `json:"interfaces"`
}

func InfoCheck(cmd Command) InfoResponse {
	info, err := host.Info()
	if err != nil {
		return InfoResponse{Error: err.Error()}
	}

	interfaces, err := net.Interfaces()
	if err != nil {
		return InfoResponse{Error: err.Error()}
	}

	resp := InfoResponse{
		Hostname: info.Hostname,
		OS:       info.OS,
		Platform: info.Platform,
	}

	for _, i := range interfaces {
		for _, up := range i.Flags {
			if up != "up" {
				continue
			}
			inter := InterfaceResponse{Name: i.Name}
			for _, addr := range i.Addrs {
				inter.IPs = append(inter.IPs, addr.Addr)
			}
			resp.Interfaces = append(resp.Interfaces, inter)
		}
	}

	return resp
}
