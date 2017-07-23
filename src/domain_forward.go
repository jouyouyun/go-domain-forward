package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func handleClientRequest(client net.Conn) {
	if client == nil {
		fmt.Println("Client nil")
		return
	}
	defer client.Close()

	var buf [1024]byte
	headerLen, err := client.Read(buf[:])
	if err != nil {
		fmt.Println("Failed to read client header:", err)
		return
	}

	method, host := parseClientHeader(buf[:])
	if host == "" {
		fmt.Println("Not found host from header:", string(buf[:]))
		return
	}

	host = convertDomain(host)
	debugLog("Host after convert: %s\n", host)

	// forward
	server, err := net.Dial("tcp", host)
	if err != nil {
		fmt.Println("Failed to dial server:", err)
		return
	}

	if method == "CONNECT" {
		fmt.Fprint(client, "HTTP/1.1 200 Connection established\r\n\r\n")
	} else {
		server.Write(buf[:headerLen])
	}
	go io.Copy(server, client)
	io.Copy(client, server)
}

func parseClientHeader(buf []byte) (method, host string) {
	var lines = strings.Split(string(buf[:]), "\n")
	for i, line := range lines {
		if line == "" {
			continue
		}

		if i == 0 {
			method = strings.Split(line, " ")[0]
			continue
		}

		if !strings.Contains(line, "Host:") {
			continue
		}

		list := strings.Split(line, " ")
		if len(list) != 2 {
			break
		}

		host = list[len(list)-1]
		if host[len(host)-1] == '\r' {
			host = host[:len(host)-1]
		}
		break
	}
	debugLog("Method: %s\nHost: %s\n", method, host)
	return
}
