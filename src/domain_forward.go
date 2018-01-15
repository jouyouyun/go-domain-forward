package main

import (
	"fmt"
	"io"
	"net"
	"strings"
	"sync"
)

func handleClientRequest(client net.Conn) {
	if client == nil {
		fmt.Println("Client nil")
		return
	}
	defer closeConnection(client)

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
	debugLog("Host after convert:", host)

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
	go func() {
		debugLog("Start copy client to server")
		_, err := doCopy(server, client)
		debugLog("Client copy to server over:", err)
		if err != nil && err != io.EOF {
			debugLog("Failed to copy in server:", err)
		}
		closeConnection(server)
	}()
	debugLog("Start copy server to client")
	_, err = doCopy(client, server)
	if err != nil {
		debugLog("Failed to copy in client:", err)
		return
	}
	debugLog("Copy done:", host)
}

var bufpool *sync.Pool

func init() {
	bufpool = &sync.Pool{}
	bufpool.New = func() interface{} {
		return make([]byte, 32*1024)
	}
}

func doCopy(dst io.Writer, src io.Reader) (written int64, err error) {
	if wt, ok := src.(io.WriterTo); ok {
		return wt.WriteTo(dst)
	}
	if rt, ok := dst.(io.ReaderFrom); ok {
		return rt.ReadFrom(src)
	}

	buf := bufpool.Get().([]byte)
	defer bufpool.Put(buf)

	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

func closeConnection(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		debugfLog("Failed to close connection: %#v\n", conn)
	}
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
	debugLog("Method and Host:", method, host)
	return
}
