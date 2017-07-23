package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

type domainInfo struct {
	Domain  string `json:"Domain"`
	Address string `json:"Address"`
}
type domainInfos []*domainInfo

var _set map[string]string

func convertDomain(domain string) string {
	domain = strings.TrimSpace(domain)
	set, err := parseConfig(*_config)
	if err != nil {
		fmt.Println("Failed to parse config:", err)
		return domain
	}

	v, ok := set[domain]
	if ok && v != "" {
		return v
	}

	tmp := strings.TrimLeft(domain, "http://")
	tmp = strings.TrimLeft(tmp, "https://")
	v, ok = set[tmp]
	if ok && v != "" {
		return v
	}
	return domain
}

func parseConfig(filename string) (map[string]string, error) {
	if _set != nil {
		return _set, nil
	}
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var infos domainInfos
	err = json.Unmarshal(content, &infos)
	if err != nil {
		return nil, err
	}

	var set = make(map[string]string)
	for _, info := range infos {
		// invalid domain or address
		if info.Domain == "" || info.Address == "" ||
			!strings.Contains(info.Address, ":") {
			continue
		}
		set[info.Domain] = info.Address
	}

	debugLog("Config file: %#v\n", set)

	_set = set
	return set, nil
}
