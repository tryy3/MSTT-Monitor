package client

import (
	"strings"

	"github.com/shirou/gopsutil/cpu"
)

// CPUInfo innehåller information om CPUn
type CPUInfo struct {
	Error     string `json:"error"`
	Cores     int32  `json:"cores"`
	ModelName string `json:"model_name"`
}

// CPUUsage innehåller information om CPU användningen
type CPUUsage struct {
	Error   string  `json:"error"`
	Procent float64 `json:"procent"`
}

// CPUCheck kollar CPUn
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

	// Kolla CPU procent från tidigare check
	// finns en chans att man får tillbaka 0 procent om
	// det är den första gången man kollar
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
