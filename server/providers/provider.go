package providers

type AlertProvider interface {
	Check(string) bool
	Value() string
	Message() string
	Name() string
}

func NewAlertProviderCPU(total int64, avg int64) AlertProvider {
	return &AlertProviderCPU{
		Total: total,
		Avg:   avg,
	}
}
