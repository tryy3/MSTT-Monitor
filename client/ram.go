package client

import (
	"strings"

	"github.com/shirou/gopsutil/mem"
)

type RAMType int

const (
	RAM RAMType = iota
	Swap
)

type RAMResponse struct {
	Error string  `json:"error"`
	Size  uint64  `json:"size"`
	Type  RAMType `json:"type"`
}

func RAMCheck(cmd Command) RAMResponse {
	t := RAM
	total := false
	resp := RAMResponse{}

	for _, args := range cmd.Params {
		if strings.ToLower(args.Name) == "-swap" {
			t = Swap
		} else if strings.ToLower(args.Name) == "-total" {
			total = true
		}
	}

	switch t {
	case RAM:
		v, err := mem.VirtualMemory()
		if err != nil {
			resp.Error = err.Error()
			return resp
		}
		if total {
			resp.Size = v.Total
		} else {
			resp.Size = v.Used
		}
		break
	case Swap:
		v, err := mem.SwapMemory()
		if err != nil {
			resp.Error = err.Error()
			return resp
		}
		if total {
			resp.Size = v.Total
		} else {
			resp.Size = v.Used
		}
		break
	}
	resp.Type = t
	return resp
}
