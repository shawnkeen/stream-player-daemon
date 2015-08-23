package main

import (
	"bufio"
	"fmt"
	"net"
)

func handleConn(conn net.Conn) {
	scanner := bufio.NewScanner(conn)
	writer := bufio.NewWriter(conn)
	writer.WriteString(protHello() + "\n")
	writer.Flush()
	for scanner.Scan() {
		line := scanner.Text()
		answer, returnCode := handleRequest(line)
		if answer != "" {
			writer.WriteString(answer + "\n")
		}
		writer.WriteString(returnCode.String() + "\n")
		writer.Flush()
	}
	fmt.Printf("connection to %s lost\n", conn.RemoteAddr())
}

func startServer() {
	listenPort := ":8080"
	listen, _ := net.Listen("tcp", listenPort)
	fmt.Println("listening on port " + listenPort + "...")
	defer listen.Close()
	for {
		conn, err := listen.Accept()
		defer conn.Close()
		if err != nil {
			continue
		}
		remoteAddr, _, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			fmt.Println(err.Error())
		}
		fmt.Printf("lookup of '%s'\n", remoteAddr)
		remoteHostName, err := net.LookupAddr(remoteAddr)
		if err != nil {
			remoteHostName = []string{remoteAddr}
			fmt.Println(err.Error())
		}
		fmt.Printf("connected to '%s'", remoteHostName[0])
		if remoteAddr != remoteHostName[0] {
			fmt.Printf(" [%s]\n", remoteAddr)
		} else {
			fmt.Println("")
		}
		go handleConn(conn)
	}
}

//conn, err := net.Listen("tcp", "localhost")
