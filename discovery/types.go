package discovery

import "go_discovery_plugin/common"

// Request defines the input for device discovery.
type Request struct {
	IpRanges    []string          `json:"ipranges"`
	Credentials []common.AuthPair `json:"credentials"`
	Port        int               `json:"port"`
}

// DeviceDiscoveryResult represents the discovery result for each device.
type DeviceDiscoveryResult struct {
	IPAddress string `json:"ipaddress"`
	Status    string `json:"status"`
	AuthPair  string `json:"credentialprofileid"`
}
type Device struct {
	IPAddress  string          `json:"ipaddress"`
	Credential common.AuthPair `json:"credential"`
	Port       int             `json:"port"`
}
