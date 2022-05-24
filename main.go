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

	portStr := os.Args[1]
	port, err := strconv.ParseInt(portStr, 10, 64)
	if err != nil {
		fmt.Printf("invalid port\n: %s", portStr)
		return
	}

	if err := run(int(port), "tcp"); err != nil {
		panic(err)
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
