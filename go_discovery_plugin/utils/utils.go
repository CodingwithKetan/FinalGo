package utils

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Ping checks if an IP address is reachable.
func Ping(ip string) bool {
	cmd := exec.Command("ping", "-c", "1", "-W", "1", ip)
	err := cmd.Run()
	if err != nil {
		log.Printf("Ping failed for %s", ip)
		return false
	}
	log.Printf("Ping successful for %s", ip)
	return true
}

// CheckPortOpen verifies if a port is open on the given IP.
func CheckPortOpen(ip string, port int) bool {
	address := fmt.Sprintf("%s:%d", ip, port)
	conn, err := net.DialTimeout("tcp", address, 2*time.Second)
	if err != nil {
		log.Printf("Port %d is closed on %s", port, ip)
		return false
	}
	conn.Close()
	log.Printf("Port %d is open on %s", port, ip)
	return true
}

// ExpandIPRanges expands IP ranges like "192.168.1.4-9" into individual IPs.
func ExpandIPRanges(ipRanges []string) []string {
	var expandedIPs []string
	for _, ipRange := range ipRanges {
		if strings.Contains(ipRange, "-") {
			expandedIPs = append(expandedIPs, expandHyphenRange(ipRange)...)
		} else {
			expandedIPs = append(expandedIPs, ipRange) // Single IP
		}
	}
	return expandedIPs
}

// expandHyphenRange converts "192.168.1.4-9" to individual IPs.
func expandHyphenRange(ipRange string) []string {
	var ips []string
	re := regexp.MustCompile(`(\d+\.\d+\.\d+)\.(\d+)-(\d+)`)
	match := re.FindStringSubmatch(ipRange)
	if len(match) == 4 {
		base := match[1] // "192.168.1"
		start, _ := strconv.Atoi(match[2])
		end, _ := strconv.Atoi(match[3])

		// Validate range
		if start > end {
			return []string{ipRange} // Invalid range, return as is
		}

		for i := start; i <= end; i++ {
			ips = append(ips, fmt.Sprintf("%s.%d", base, i))
		}
	} else {
		ips = append(ips, ipRange) // If not a range, return as is
	}
	return ips
}
