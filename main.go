package main

import (
	"fmt"
	"net"
	"sort"
)

func worker(ports chan int, results chan int) {
	for p := range ports {
		address := fmt.Sprintf("20.194.168.9:%d", p)
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
	ports := make(chan int, 100)
	results := make(chan int)
	var openPorts []int
	var closedPorts []int

	for i := 0; i < cap(ports); i++ {
		go worker(ports, results)
	}

	go func() {
		for i := 0; i < 1024; i++ {
			ports <- i
		}
	}()

	for i := 0; i < 1024; i++ {
		port := <-results
		if port != 0 {
			openPorts = append(openPorts, port)
		} else {
			closedPorts = append(closedPorts, port)
		}
	}

	close(ports)
	close(results)

	// all closed ports are 0
	// sort.Ints(closedPorts)
	// for _, port := range closedPorts {
	// 	fmt.Printf("%d is closed\n", port)
	// }

	sort.Ints(openPorts)
	for _, port := range openPorts {
		fmt.Printf("%d is open\n", port)
	}
}
