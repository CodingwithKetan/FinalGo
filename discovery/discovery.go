package discovery

import (
	"log"
	"sync"

	"go_discovery_plugin/common"
	"go_discovery_plugin/sshutils"
	"go_discovery_plugin/utils"
)

// RunDiscovery executes the discovery process with one goroutine per IP.
func RunDiscovery(request Request) []DeviceDiscoveryResult {
	var results []DeviceDiscoveryResult
	ipList := utils.ExpandIPRanges(request.IpRanges)

	var wg sync.WaitGroup
	resultsChan := make(chan DeviceDiscoveryResult, len(ipList)) // Channel to collect results

	// Spawn one goroutine per IP address
	for _, ip := range ipList {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()
			resultsChan <- discoverDevice(ip, request.Port, request.Credentials)
		}(ip)
	}

	// Wait for all goroutines to finish
	go func() {
		wg.Wait()
		close(resultsChan) // Close channel after all workers are done
	}()

	// Collect results from channel
	for res := range resultsChan {
		results = append(results, res)
	}

	return results
}

// discoverDevice handles discovery for a single device (ping, port, SSH).
func discoverDevice(ip string, port int, credentials []common.AuthPair) DeviceDiscoveryResult {
	// Step 1: Ping Check
	if !utils.Ping(ip) {
		log.Printf("Ping failed for %s", ip)
		return DeviceDiscoveryResult{
			IPAddress: ip,
			Status:    "Ping Failed",
			AuthPair:  "N/A",
		}
	}

	// Step 2: Port Check
	if !utils.CheckPortOpen(ip, port) {
		log.Printf("Port %d closed for %s", port, ip)
		return DeviceDiscoveryResult{
			IPAddress: ip,
			Status:    "Port Closed",
			AuthPair:  "N/A",
		}
	}

	// Step 3: SSH Authentication and Command Execution
	for _, cred := range credentials {
		if result := attemptSSH(ip, port, cred); result.Status == "Succeed" {
			return result // Return immediately on first successful authentication
		}
	}

	// If no credentials worked, mark as non-discoverable
	return DeviceDiscoveryResult{
		IPAddress: ip,
		Status:    "Non-Discoverable",
		AuthPair:  "N/A",
	}
}

// attemptSSH tries SSH authentication and runs a command.
func attemptSSH(ip string, port int, cred common.AuthPair) DeviceDiscoveryResult {
	log.Printf("Trying SSH authentication for %s with user %s", ip, cred.Username)

	conn, success := sshutils.SSHAuth(ip, port, cred.Username, cred.Password)
	if success {
		defer conn.Close()
		output := sshutils.RunCommand(conn, "whoami")
		if output != "" {
			log.Printf("Device %s successfully discovered", ip)
			return DeviceDiscoveryResult{
				IPAddress: ip,
				Status:    "Succeed",
				AuthPair:  cred.ID,
			}
		} else {
			log.Printf("Device %s connected but unable to execute command", ip)
		}
	}

	return DeviceDiscoveryResult{
		IPAddress: ip,
		Status:    "SSH Failed",
		AuthPair:  cred.ID,
	}
}
