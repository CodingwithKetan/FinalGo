package discovery

import (
	"log"
	"sync"
	"time"

	"go_discovery_plugin/sshutils"
	"go_discovery_plugin/utils"
)

// RunDiscovery executes the discovery process using channels.
func RunDiscovery(request Request) []DeviceDiscoveryResult {
	var results []DeviceDiscoveryResult
	ipList := utils.ExpandIPRanges(request.IpRanges)

	var wg sync.WaitGroup
	resultsChan := make(chan DeviceDiscoveryResult, len(ipList)) // Channel to collect results

	// Step 1: Ping and Port Check
	for _, ip := range ipList {
		if !utils.Ping(ip) {
			resultsChan <- DeviceDiscoveryResult{
				IPAddress:       ip,
				Status:          "Ping Failed",
				AuthPair:        "N/A",
				LastCheckedTime: time.Now().Format(time.RFC3339),
			}
			continue
		}

		if !utils.CheckPortOpen(ip, request.Port) {
			resultsChan <- DeviceDiscoveryResult{
				IPAddress:       ip,
				Status:          "Port Closed",
				AuthPair:        "N/A",
				LastCheckedTime: time.Now().Format(time.RFC3339),
			}
			continue
		}

		// Step 2: Try SSH Authentication (Parallel Execution)
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()

			var result DeviceDiscoveryResult
			for _, cred := range request.Credentials {
				result = processDiscovery(Device{
					IPAddress: ip, Credential: cred, Port: request.Port,
				})

				// If SSH authentication succeeds, stop checking further credentials
				if result.Status == "Discoverable" {
					break
				}
			}

			// Send the result to the channel
			resultsChan <- result
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

// processDiscovery attempts SSH authentication for a device.
func processDiscovery(device Device) DeviceDiscoveryResult {
	log.Printf("Trying SSH authentication for %s", device.IPAddress)

	conn, success := sshutils.SSHAuth(device.IPAddress, device.Port, device.Credential.Username, device.Credential.Password)
	if success {
		defer conn.Close()
		output := sshutils.RunCommand(conn, "whoami")
		if output != "" {
			log.Printf("Device %s successfully discovered", device.IPAddress)
			return DeviceDiscoveryResult{
				IPAddress:       device.IPAddress,
				Status:          "Discoverable",
				AuthPair:        device.Credential.ID,
				LastCheckedTime: time.Now().Format(time.RFC3339),
			}
		} else {
			log.Printf("Device %s connected but unable to execute command", device.IPAddress)
			return DeviceDiscoveryResult{
				IPAddress:       device.IPAddress,
				Status:          "Connected but Command Failed",
				AuthPair:        device.Credential.ID,
				LastCheckedTime: time.Now().Format(time.RFC3339),
			}
		}
	}

	return DeviceDiscoveryResult{
		IPAddress:       device.IPAddress,
		Status:          "SSH Authentication Failed",
		AuthPair:        "N/A",
		LastCheckedTime: time.Now().Format(time.RFC3339),
	}
}
