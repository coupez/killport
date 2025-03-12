package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <port1> [port2] [port3] ...\n", os.Args[0])
		os.Exit(1)
	}

	failed := []string{}
	portStrings := os.Args[1:]

	for _, portStr := range portStrings {
		port, err := strconv.ParseInt(portStr, 10, 64)
		if err != nil {
			fmt.Printf("Invalid port: %s\n", portStr)
			failed = append(failed, portStr)
			continue
		}

		// Check port range
		if port < 1 || port > 65535 {
			fmt.Printf("Port out of range (1-65535): %d\n", port)
			failed = append(failed, portStr)
			continue
		}

		if err := killProcessOnPort(int(port), "tcp"); err != nil {
			fmt.Printf("Failed to kill process on TCP port %d: %v\n", port, err)
			failed = append(failed, portStr)
		} else {
			fmt.Printf("Successfully killed process on TCP port %d\n", port)
		}
	}

	if len(failed) > 0 {
		fmt.Printf("Failed ports: %s\n", strings.Join(failed, ", "))
		os.Exit(1)
	}
}

func killProcessOnPort(port int, protocol string) error {
	os := runtime.GOOS

	switch os {
	case "windows":
		// For Windows using netstat and taskkill
		findCmd := fmt.Sprintf("netstat -ano | findstr :%d | findstr %s", port, strings.ToUpper(protocol))
		findOutput, err := exec.Command("cmd", "/C", findCmd).Output()
		if err != nil {
			// No process found on this port, which is not an error
			return nil
		}

		// Parse the output to get PIDs
		lines := strings.Split(string(findOutput), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) < 5 {
				continue
			}
			pid := fields[4]
			killCmd := fmt.Sprintf("taskkill /F /PID %s", pid)
			_, err = exec.Command("cmd", "/C", killCmd).Output()
			if err != nil {
				return fmt.Errorf("failed to kill process with PID %s: %v", pid, err)
			}
		}

	case "darwin", "linux":
		// For macOS and Linux using lsof
		var grepPattern string
		if protocol == "tcp" {
			grepPattern = "LISTEN"
		} else {
			grepPattern = "UDP"
		}

		findCmd := fmt.Sprintf("lsof -i %s:%d | grep %s", protocol, port, grepPattern)
		findOutput, err := exec.Command("sh", "-c", findCmd).Output()
		if err != nil {
			// No process found on this port, which is not an error
			return nil
		}

		// Parse the output to get PIDs
		lines := strings.Split(string(findOutput), "\n")
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			fields := strings.Fields(line)
			if len(fields) < 2 {
				continue
			}
			pid := fields[1]
			killCmd := fmt.Sprintf("kill -9 %s", pid)
			_, err = exec.Command("sh", "-c", killCmd).Output()
			if err != nil {
				return fmt.Errorf("failed to kill process with PID %s: %v", pid, err)
			}
		}

	default:
		return fmt.Errorf("unsupported operating system: %s", os)
	}

	return nil
}
