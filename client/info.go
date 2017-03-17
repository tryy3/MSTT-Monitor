package client

import (
	"strings"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/net"
)

// UptimeResponse innehåller information om uptime på klienten
type UptimeResponse struct {
	Error  string `json:"error"`
	Uptime uint64 `json:"uptime"`
}

// UptimeCheck kollar hur länge en klient har varit aktiv
func UptimeCheck(cmd Command) UptimeResponse {
	boot := false
	for _, args := range cmd.Params {
		if strings.ToLower(args.Name) == "-boot" {
			boot = true
		}
	}

	if boot {
		b, err := host.BootTime()
		if err != nil {
			return UptimeResponse{Error: err.Error()}
		}
		return UptimeResponse{Uptime: b}
	}

	up, err := host.Uptime()
	if err != nil {
		return UptimeResponse{Error: err.Error()}
	}
	return UptimeResponse{Uptime: up}
}

// InterfaceResponse innehåller information om en nätverks interface
type InterfaceResponse struct {
	Name string   `json:"name"`
	IPs  []string `json:"ips"` // Kan ha både ipv4 och ipv6
}

// InfoResponse innehåller information om klienten
type InfoResponse struct {
	Error         string              `json:"error"`
	Hostname      string              `json:"hostname"`
	OS            string              `json:"os"`
	Platform      string              `json:"platform"`
	ClientVersion string              `json:"client_version"`
	Interfaces    []InterfaceResponse `json:"interfaces"`
}

// InfoCheck hämtar information om klienten
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
		Hostname:      info.Hostname,
		OS:            info.OS,
		Platform:      info.Platform,
		ClientVersion: version,
	}

	for _, i := range interfaces {
		for _, up := range i.Flags {
			// Kolla om nätverks interfacen är aktiv
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
