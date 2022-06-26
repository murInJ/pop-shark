package pop_shark

import (
	"encoding/json"
	"errors"
	"net"
)

func jsonStr2map(str string) (map[string]interface{}, error) {
	var t interface{}
	err := json.Unmarshal([]byte(str), &t)

	if err != nil {
		return nil, err
	}

	return t.(map[string]interface{}), nil
}

func data2jsonStr(data interface{}) (string, error) {
	b, err := json.Marshal(data)

	if err != nil {
		return "", err
	}

	return string(b), nil
}

func externalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			ip := getIpFromAddr(addr)
			if ip == nil {
				continue
			}
			return ip, nil
		}
	}
	return nil, errors.New("connected to the network?")
}

func getIpFromAddr(addr net.Addr) net.IP {
	var ip net.IP
	switch v := addr.(type) {
	case *net.IPNet:
		ip = v.IP
	case *net.IPAddr:
		ip = v.IP
	}
	if ip == nil || ip.IsLoopback() {
		return nil
	}
	ip = ip.To4()
	if ip == nil {
		return nil // not an ipv4 address
	}

	return ip
}
