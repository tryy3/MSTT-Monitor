package client

import (
	"strings"

	"github.com/shirou/gopsutil/disk"
)

type PartitionResponse struct {
	Name string `json:"name"`
	Size uint64 `json:"size"`
	Type string `json:"type"`
}

type DiscResponse struct {
	Error      string              `json:"error"`
	Partitions []PartitionResponse `json:"partitions"`
}

func DiscCheck(cmd Command) DiscResponse {
	part, err := disk.Partitions(false)
	if err != nil {
		return DiscResponse{Error: err.Error()}
	}

	resp := DiscResponse{}
	total := false

	for _, args := range cmd.Params {
		if strings.ToLower(args.Name) == "-total" {
			total = true
		}
	}

	for _, p := range part {
		usage, err := disk.Usage(p.Mountpoint)
		if err != nil {
			return DiscResponse{Error: err.Error()}
		}
		partition := PartitionResponse{
			Name: p.Mountpoint,
			Type: p.Fstype,
		}

		if total {
			partition.Size = usage.Total
		} else {
			partition.Size = usage.Used
		}

		resp.Partitions = append(resp.Partitions, partition)
	}
	return resp
}
