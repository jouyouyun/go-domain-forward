package main

import (
	"fmt"
	"testing"
)

func TestConfig(t *testing.T) {
	src := *_config
	*_config = "sample.json"
	var infos = domainInfos{
		{
			Domain:  "blog.test.org",
			Address: "127.0.0.1:8000",
		},
		{
			Domain:  "http://note.test.org",
			Address: "127.0.0.1:9000",
		},
		{
			Domain:  "https://wiki.test.org",
			Address: "127.0.0.1:8080",
		},
		{
			Domain:  "test.org",
			Address: "test.org",
		},
		{
			Domain:  "http://test.org",
			Address: "http://test.org",
		},
		{
			Domain:  "https://test.org",
			Address: "https://test.org",
		},
	}
	for _, info := range infos {
		v := convertDomain(info.Domain)
		if v != info.Address {
			panic(fmt.Sprintf("Expected: %q\nBut obtained: %q\n", info.Address, v))
		}
	}
	*_config = src
}

func TestParseHeader(t *testing.T) {
	var header = `GET /favicon.ico HTTP/1.1
Host: wiki.test.org
User-Agent: Mozilla/5.0 (X11; Linux x86_64; rv:54.0) Gecko/20100101 Firefox/54.0
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
Accept-Language: en-US,en;q=0.5
Accept-Encoding: gzip, deflate
Connection: keep-alive

`
	method, host := parseClientHeader([]byte(header))
	if method != "GET" {
		panic(fmt.Sprintf("Expected: \"GET\"\nBut obtained: %q\n", method))
	}
	if host != "wiki.test.org" {
		panic(fmt.Sprintf("Expected: \"GET\"\nBut obtained: %q\n", host))
	}
}
