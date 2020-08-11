package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
)

var (
	successCount int
	successMax   *int
	startTime    int64
	port         *int
)

func TCPPing(IP string, port *int, waitTimeout time.Duration) {
	conn, err := net.DialTimeout("tcp", IP+":"+strconv.Itoa(*port), waitTimeout*time.Second)
	if err != nil {
		//fmt.Printf("\r\033[2KPort %v closed (worker: %v)", i, workerID)
	} else {
		conn.Close()
		successCount++
		output := fmt.Sprintf("Port %v open", *port)
		fmt.Print(output + "\n")

		output = "The Internet's on!"
		cmd := exec.Command("say", "-v", "Good News", output)
		cmd.Run()

		if successCount < *successMax {
			time.Sleep(waitTimeout * time.Second)
		}
	}
}

func IPFormat(IP string) string {
	for i := 0; i < len(IP); i++ {
		switch IP[i] {
		case '.':
			return ("IPv4")
		case ':':
			return ("IPv6")
		}
	}
	return ""
}

func checkValidIP(IP string) {
	if net.ParseIP(IP) == nil {
		fmt.Printf("Invalid IP Address: %s\n", IP)
		os.Exit(1)
	}
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func handleSIGTERM() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Printf("\n\nPing cancelled\n")
		stats()
		os.Exit(0)
	}()
}

func stats() {
	td := time.Now()
	endTime := td.Unix()
	timeTaken := endTime - startTime
	fmt.Printf("\r\nPinged port %v for %v seconds.\n", *port, timeTaken)
}

func main() {
	handleSIGTERM()

	usage := fmt.Sprintf("Usage: %s -p <port> -s <stop after X successful pings>-w <wait timeout> <ip address>\n", filepath.Base(os.Args[0]))

	argLength := len(os.Args[1:])
	if argLength < 1 {
		fmt.Printf(usage)
		os.Exit(1)
	}

	//set defaults
	port = flag.Int("p", 80, "start port")
	successMax = flag.Int("s", 3, "stop after X successful pings")
	waitTimeout := flag.Int64("w", 30, "Wait timeout")

	flag.Parse()
	IP := flag.Arg(0)

	checkValidIP(IP)

	if *port < 0 || *port > 65535 {
		*port = 80
	}

	successCount = 0
	fmt.Printf("TCP Ping %s (%s) Port %v\n\n", IP, IPFormat(IP), *port)
	td := time.Now()
	startTime = td.Unix()

	for successCount < *successMax {
		TCPPing(IP, port, time.Duration(*waitTimeout))
	}

	stats()
}
