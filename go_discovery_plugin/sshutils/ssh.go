package sshutils

import (
	"fmt"
	"log"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
)

// SSHAuth verifies SSH authentication using username & password.
func SSHAuth(ip string, port int, username, password string) (*ssh.Client, bool) {
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         3 * time.Second,
	}

	addr := fmt.Sprintf("%s:%d", ip, port)
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		log.Printf("SSH Authentication failed for %s (%s)", ip, username)
		return nil, false
	}
	log.Printf("SSH Authentication successful for %s (%s)", ip, username)
	return conn, true
}

// RunCommand executes a command on the remote machine via SSH.
func RunCommand(conn *ssh.Client, cmd string) string {
	session, err := conn.NewSession()
	if err != nil {
		log.Printf("Failed to create SSH session: %v", err)
		return "N/A"
	}
	defer session.Close()

	output, err := session.CombinedOutput(cmd)
	if err != nil {
		log.Printf("Error executing command '%s': %v", cmd, err)
		return "N/A"
	}

	return strings.TrimSpace(string(output))
}

// FetchMetrics executes system monitoring commands on a Linux server.
func FetchMetrics(conn *ssh.Client) map[string]string {
	metrics := make(map[string]string)

	metrics["cpu"] = RunCommand(conn, `top -bn1 | grep 'Cpu(s)' | awk '{print $2}'`)
	metrics["memory"] = RunCommand(conn, `free -m | awk 'NR==2{printf "%.2f%%", $3*100/$2 }'`)
	metrics["disk"] = RunCommand(conn, `df -h / | awk 'NR==2 {print $5}'`)

	return metrics
}
