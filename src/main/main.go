package main

import (
	"fmt"
	"net"
	"sort"
)

const (
	ports_count         = 1024
	channel_buffer_size = 100
	url                 = "scanme.nmap.org"
)

func worker(ports, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("%s:%d", url, p)
		conn, err := net.Dial("tcp", address)
		if err != nil {
			results <- 0
			continue
		}
		conn.Close()
		results <- p
	}
}

func main() {
	ports := make(chan int, channel_buffer_size)
	results := make(chan int)
	var openports []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 1; i <= ports_count; i++ {
			ports <- i
		}
	}()

	for i := 1; i <= ports_count; i++ {
		port := <-results
		if port != 0 {
			openports = append(openports, port)
		}
	}

	close(ports)
	close(results)
	sort.Ints(openports)
	for _, port := range openports {
		fmt.Printf("%d open\n", port)
	}
}
