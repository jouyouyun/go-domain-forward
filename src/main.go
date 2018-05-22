package main

import (
	"flag"
	"fmt"
	"net"
)

var (
	_port   = flag.String("p", "8000", "port")
	_config = flag.String("c", "/etc/go-domain-forward/config.json", "config")
	_debug  = flag.Bool("d", false, "debug")
)

func debugLog(format string, a ...interface{}) {
	if *_debug == false {
		return
	}
	fmt.Println(format, a)
}

func debugfLog(format string, a ...interface{}) {
	if *_debug == false {
		return
	}
	fmt.Printf(format, a)
}

func main() {
	flag.Parse()
	l, err := net.Listen("tcp", ":"+(*_port))
	if err != nil {
		fmt.Println("Failed to listen:", err)
		return
	}

	fmt.Println("Wait connecting...")
	for {
		request, err := l.Accept()
		if err != nil {
			fmt.Println("Failed to accept:", err)
			continue
		}

		go handleClientRequest(request)
	}
}
