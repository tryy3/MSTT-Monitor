package providers

import (
	"sync"
)

type AlertProvider interface {
	Check(string) bool
	Value() string
	Message() string
	Name() string
}

func NewAlertProviderCPU(total int64, avg float64) AlertProvider {
	return &AlertProviderCPU{
		rw:    new(sync.RWMutex),
		Total: total,
		Avg:   avg,
	}
}
