package metrics

import "go_discovery_plugin/common"

// Request defines the input for metrics collection.
type Request struct {
	IpRanges    []string          `json:"ipranges"`
	Credentials []common.AuthPair `json:"credentials"`
	Port        int               `json:"port"`
}

// Device represents a device with authentication credentials.
type Device struct {
	IPAddress string `json:"ipaddress"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Port      int    `json:"port"`
}

// MetricsResult contains system metrics for a device.
type MetricsResult struct {
	IPAddress   string `json:"ipaddress"`
	CPUUsage    string `json:"cpuusage"`
	MemoryUsage string `json:"memoryusage"`
	DiskUsage   string `json:"diskusage"`
}
