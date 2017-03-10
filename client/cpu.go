package client

import (
	"strings"

	"github.com/shirou/gopsutil/cpu"
)

type CPUInfo struct {
	Error     string `json:"error"`
	Cores     int32  `json:"cores"`
	ModelName string `json:"model_name"`
}

type CPUUsage struct {
	Error   string  `json:"error"`
	Procent float64 `json:"procent"`
}

func CPUCheck(cmd Command) interface{} {
	info := false

	for _, args := range cmd.Params {
		if strings.ToLower(args.Name) == "-info" {
			info = true
		}
	}

	if info {
		info, err := cpu.Info()
		if err != nil {
			return CPUInfo{Error: err.Error()}
		}

		if len(info) <= 0 {
			return CPUInfo{Error: "Can't find CPU info"}
		}

		i := info[0]
		return CPUInfo{
			Cores:     i.Cores,
			ModelName: i.ModelName,
		}
	}

	f, err := cpu.Percent(0, false)
	if err != nil {
		return CPUUsage{Error: err.Error()}
	}

	if len(f) <= 0 {
		return CPUInfo{Error: "Can't find CPU usage"}
	}

	i := f[0]

	return CPUUsage{Procent: i}
}
