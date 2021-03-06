package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("missing port\n")
		return
	}

	failed := []int64{}
	portStrings := os.Args[1:]

	for _, portStr := range portStrings {
		port, err := strconv.ParseInt(portStr, 10, 64)
		if err != nil {
			fmt.Printf("invalid port: %s\n", portStr)
			failed = append(failed, port)
		}

		if err := run(int(port), "tcp"); err != nil {
			panic(err)
		}
	}

	if len(failed) != 0 {
		os.Exit(1)
	}
}

func run(port int, protocol string) error {
	grepLine := "LISTEN"
	if protocol == "udp" {
		grepLine = "UDP"
	}

	killCmd := fmt.Sprintf("lsof -i %s:%d | grep %s | awk '{print $2}' | xargs kill", protocol, port, grepLine)
	cmd := exec.Command("sh", "-c", killCmd)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Println(string(stdout))
		return err
	}
	return nil
}
