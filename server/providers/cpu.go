package providers

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// CPUUsage innehåller information om CPU användningen
type CPUUsage struct {
	Error   string  `json:"error"`
	Procent float64 `json:"procent"`
}

type AlertProviderCPU struct {
	Values []float64
	Avg    int64
	Total  int64
}

func (a AlertProviderCPU) Name() string {
	return "CPU"
}

func (a *AlertProviderCPU) Check(resp string) bool {
	usage := CPUUsage{}
	err := json.Unmarshal([]byte(resp), &usage)
	if err != nil || usage.Error != "" {
		return false
	}

	if int64(len(a.Values)) >= a.Total {
		a.Values = append(a.Values[1:], usage.Procent)
	} else {
		a.Values = append(a.Values, usage.Procent)
	}

	if a.avg() > float64(a.Avg) {
		return true
	}
	return false
}

func (a *AlertProviderCPU) avg() float64 {
	var i float64 = 0
	for _, v := range a.Values {
		i += v
	}
	return i / float64(len(a.Values))
}

func (a AlertProviderCPU) Value() string {
	return strconv.FormatFloat(a.avg(), 'g', -1, 64)
}

func (a AlertProviderCPU) Message() string {
	return fmt.Sprintf("CPU Usage: %s", a.Value())
}
