package client

import (
	"strings"

	"github.com/shirou/gopsutil/mem"
)

// MemoryType är en typ för vilken typ
// av minne som ska användas
type MemoryType int

const (
	RAM  MemoryType = iota // 0
	Swap                   // 1
)

// MemoryResponse innehåller information om RAM/Swap
type MemoryResponse struct {
	Error string     `json:"error"`
	Size  uint64     `json:"size"`
	Type  MemoryType `json:"type"`
}

// MemoryCheck kollar användningen av ett minnestyp
func MemoryCheck(cmd Command) MemoryResponse {
	t := RAM
	total := false
	resp := MemoryResponse{Type: t}

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

	return resp
}
