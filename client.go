package main

import (
	"bufio"
	"fmt"
	"net"
)

func client() {
	conn, _ := net.Dial("tcp", "localhost:8080")
	defer conn.Close()
	scanner := bufio.NewScanner(conn)
	answer := new([]string)
	for scanner.Scan() {
		line := scanner.Text()
		if protResponseEnd(line) {

		}
		answer = append(answer, line)
		fmt.Println(ret)
	}
}
