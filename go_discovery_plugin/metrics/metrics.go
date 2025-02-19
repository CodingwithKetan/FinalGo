package metrics

import (
	"log"
	"sync"

	"go_discovery_plugin/sshutils"
	"go_discovery_plugin/utils"
)

// RunMetricsCollection runs the metrics collection process using channels.
func RunMetricsCollection(request Request) []MetricsResult {
	ipList := utils.ExpandIPRanges(request.IpRanges)

	var wg sync.WaitGroup
	resultsChan := make(chan MetricsResult, len(ipList)) // Buffered channel to store results (one per IP)

	// Step 1: Ping Check (Before SSH)
	for _, ip := range ipList {
		wg.Add(1)
		go func(ip string) {
			defer wg.Done()

			// If IP is not reachable, return "N/A"
			if !utils.Ping(ip) {
				resultsChan <- MetricsResult{
					IPAddress:   ip,
					CPUUsage:    "N/A",
					MemoryUsage: "N/A",
					DiskUsage:   "N/A",
				}
				return
			}

			// Step 2: Try multiple credentials sequentially
			var result MetricsResult
			success := false
			for _, cred := range request.Credentials {
				result = collectMetrics(Device{
					IPAddress: ip, Username: cred.Username, Password: cred.Password, Port: request.Port,
				})

				// If SSH connection is successful, stop trying further credentials
				if result.CPUUsage != "N/A" {
					success = true
					break
				}
			}

			// If no credential worked, return "N/A"
			if !success {
				result = MetricsResult{
					IPAddress:   ip,
					CPUUsage:    "N/A",
					MemoryUsage: "N/A",
					DiskUsage:   "N/A",
				}

			}

			resultsChan <- result // Send the final result to the channel
		}(ip)
	}

	// Step 3: Close channel after all IPs are processed
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// Step 4: Collect results from channel
	var results []MetricsResult
	for result := range resultsChan {
		results = append(results, result)
	}

	return results
}

// collectMetrics connects via SSH and retrieves metrics.
func collectMetrics(device Device) MetricsResult {
	log.Printf("Fetching metrics for %s", device.IPAddress)

	conn, success := sshutils.SSHAuth(device.IPAddress, device.Port, device.Username, device.Password)
	if !success {
		return MetricsResult{
			IPAddress:   device.IPAddress,
			CPUUsage:    "N/A",
			MemoryUsage: "N/A",
			DiskUsage:   "N/A",
		}
	}
	defer conn.Close()

	metrics := sshutils.FetchMetrics(conn)
	return MetricsResult{
		IPAddress:   device.IPAddress,
		CPUUsage:    metrics["cpu"],
		MemoryUsage: metrics["memory"],
		DiskUsage:   metrics["disk"],
	}
}
