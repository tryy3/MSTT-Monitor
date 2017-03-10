package server

import (
	"encoding/json"
	"testing"
	"time"
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

type PartitionResponse struct {
	Name string `json:"name"`
	Size uint64 `json:"size"`
	Type string `json:"type"`
}

type DiscResponse struct {
	Error      string              `json:"error"`
	Partitions []PartitionResponse `json:"partitions"`
}

type File struct {
	Name    string    `json:"name"`
	Size    int64     `json:"size"`
	IsDir   bool      `json:"isdir"`
	Perms   string    `json:"perms"`
	LastMod time.Time `json:"lastmod"`
}

type FileResponse struct {
	Error string `json:"error"`
	MFile File   `json:"mfile"`
}

type UptimeResponse struct {
	Error  string `json:"error"`
	Uptime uint64 `json:"uptime"`
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

type RAMResponse struct {
	Error string `json:"error"`
	Size  uint64 `json:"size"`
	Type  int    `json:"type"`
}

func TestRam(t *testing.T) {
	resp, err := SendMessage("192.168.20.164", "3333", "tcp", "check_ram")
	if err != nil {
		t.Error(err)
	}

	r := &RAMResponse{}
	err = json.Unmarshal([]byte(resp), r)
	if err != nil {
		t.Error(err)
	}

	if r.Error != "" {
		t.Error(r.Error)
	}

	if r.Type != 0 {
		t.Error("For check_ram", "Type", "expected", 0, "got", r.Type)
	}

	if r.Size <= 0 {
		t.Error("For check_ram", "Size", "expected a number", "got", r.Size)
	}
}

func TestRamTotal(t *testing.T) {
	resp, err := SendMessage("192.168.20.164", "3333", "tcp", "check_ram -total")
	if err != nil {
		t.Error(err)
	}

	r := &RAMResponse{}
	err = json.Unmarshal([]byte(resp), r)
	if err != nil {
		t.Error(err)
	}

	if r.Error != "" {
		t.Error(r.Error)
	}

	if r.Type != 0 {
		t.Error("For check_ram -total", "Type", "expected", 0, "got", r.Type)
	}

	if r.Size <= 0 {
		t.Error("For check_ram -total", "Size", "expected a number", "got", r.Size)
	}
}

func TestSwap(t *testing.T) {
	resp, err := SendMessage("192.168.20.164", "3333", "tcp", "check_ram -swap")
	if err != nil {
		t.Error(err)
	}

	r := &RAMResponse{}
	err = json.Unmarshal([]byte(resp), r)
	if err != nil {
		t.Error(err)
	}

	if r.Error != "" {
		t.Error(r.Error)
	}

	if r.Type != 1 {
		t.Error("For check_ram -swap", "Type", "expected", 1, "got", r.Type)
	}
}

func TestSwapTotal(t *testing.T) {
	resp, err := SendMessage("192.168.20.164", "3333", "tcp", "check_ram -swap -total")
	if err != nil {
		t.Error(err)
	}

	r := &RAMResponse{}
	err = json.Unmarshal([]byte(resp), r)
	if err != nil {
		t.Error(err)
	}

	if r.Error != "" {
		t.Error(r.Error)
	}

	if r.Type != 1 {
		t.Error("For check_ram -swap -total", "Type", "expected", 1, "got", r.Type)
	}
}

func TestDisc(t *testing.T) {
	resp, err := SendMessage("192.168.20.164", "3333", "tcp", "check_disc")
	if err != nil {
		t.Error(err)
	}

	r := &DiscResponse{}
	err = json.Unmarshal([]byte(resp), r)
	if err != nil {
		t.Error(err)
	}

	if r.Error != "" {
		t.Error(r.Error)
	}

	if len(r.Partitions) <= 0 {
		t.Error("For check_disc", "Partition length", "expected a number greater then", "0", "got", len(r.Partitions))
	}

	for _, parts := range r.Partitions {
		if parts.Size <= 0 {
			t.Error("For check_disc", "Partition size", "expected a number greated then", "0", "got", parts.Size)
		}

		if parts.Name == "" {
			t.Error("For check_disc", "Partition name", "expected a string", "got", parts.Name)
		}

		if parts.Type == "" {
			t.Error("For check_disc", "Partition type", "expected a string", "got", parts.Type)
		}
	}
}

func TestDiscTotal(t *testing.T) {
	resp, err := SendMessage("192.168.20.164", "3333", "tcp", "check_disc -total")
	if err != nil {
		t.Error(err)
	}

	r := &DiscResponse{}
	err = json.Unmarshal([]byte(resp), r)
	if err != nil {
		t.Error(err)
	}

	if r.Error != "" {
		t.Error(r.Error)
	}

	if len(r.Partitions) <= 0 {
		t.Error("For check_disc -total", "Partition size", "expected a number greater then 0", "got", len(r.Partitions))
	}

	for _, parts := range r.Partitions {
		if parts.Size <= 0 {
			t.Error("For check_disc -total", "Partition size", "expected a number greated then", "0", "got", parts.Size)
		}

		if parts.Name == "" {
			t.Error("For check_disc -total", "Partition name", "expected a string", "got", parts.Name)
		}

		if parts.Type == "" {
			t.Error("For check_disc -total", "Partition type", "expected a string", "got", parts.Type)
		}
	}
}

func TestCPU(t *testing.T) {
	resp, err := SendMessage("192.168.20.164", "3333", "tcp", "check_cpu")
	if err != nil {
		t.Error(err)
	}

	r := &CPUUsage{}
	err = json.Unmarshal([]byte(resp), r)
	if err != nil {
		t.Error(err)
	}

	if r.Error != "" {
		t.Error(r.Error)
	}

	if r.Procent <= 0 {
		t.Error("For check_cpu", "CPU Procent", "expected a number greated then", "0", "got", r.Procent)
	}
}

func TestCPUInfo(t *testing.T) {
	resp, err := SendMessage("192.168.20.164", "3333", "tcp", "check_cpu -info")
	if err != nil {
		t.Error(err)
	}

	r := &CPUInfo{}
	err = json.Unmarshal([]byte(resp), r)
	if err != nil {
		t.Error(err)
	}

	if r.Error != "" {
		t.Error(r.Error)
	}

	if r.Cores <= 0 {
		t.Error("For check_cpu -info", "CPU Cores", "expected a number greated then", "0", "got", r.Cores)
	}

	if r.ModelName == "" {
		t.Error("For check_cpu -info", "CPU Model", "expected a string", "got", r.ModelName)
	}
}

func TestUptime(t *testing.T) {
	resp, err := SendMessage("192.168.20.164", "3333", "tcp", "uptime")
	if err != nil {
		t.Error(err)
	}

	r := &UptimeResponse{}
	err = json.Unmarshal([]byte(resp), r)
	if err != nil {
		t.Error(err)
	}

	if r.Error != "" {
		t.Error(r.Error)
	}

	if r.Uptime <= 0 {
		t.Error("For uptime", "expected a number greated then", "0", "got", r.Uptime)
	}
}

func TestInfo(t *testing.T) {
	resp, err := SendMessage("192.168.20.164", "3333", "tcp", "info")
	if err != nil {
		t.Error(err)
	}

	r := &InfoResponse{}
	err = json.Unmarshal([]byte(resp), r)
	if err != nil {
		t.Error(err)
	}

	if r.Error != "" {
		t.Error(r.Error)
	}

	if r.Hostname == "" {
		t.Error("For info", "Hostname", "expected a string", "got", r.Hostname)
	}

	if r.OS == "" {
		t.Error("For info", "OS", "expected a string", "got", r.OS)
	}

	if r.Platform == "" {
		t.Error("For info", "Platform", "expected a string", "got", r.Platform)
	}
}

func TestFile(t *testing.T) {
	resp, err := SendMessage("192.168.20.164", "3333", "tcp", "file -file=C:/Users/Dennis/Downloads")
	if err != nil {
		t.Error(err)
	}

	r := &FileResponse{}
	err = json.Unmarshal([]byte(resp), r)
	if err != nil {
		t.Error(err)
	}

	if r.Error != "" {
		t.Error(r.Error)
	}
}
