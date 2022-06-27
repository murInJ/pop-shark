package pop_shark

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/murInJ/amazonsChess"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
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

func toIntSlice(actual interface{}) ([]int, error) {
	var res []int
	value := reflect.ValueOf(actual)
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return nil, errors.New("parse error")
	}
	eleType := value.Index(0).Kind()
	if eleType == reflect.Float64 {
		for i := 0; i < value.Len(); i++ {
			res = append(res, int(value.Index(i).Interface().(float64)))
		}
	} else if eleType == reflect.Int {
		for i := 0; i < value.Len(); i++ {
			res = append(res, value.Index(i).Interface().(int))
		}
	} else if eleType == reflect.Interface {
		for i := 0; i < value.Len(); i++ {
			res = append(res, int(value.Index(i).Interface().(float64)))
		}
	}

	return res, nil
}

func getIp() (string, error) {
	responseClient, errClient := http.Get("http://ip.dhcp.cn/?ip") // 获取外网 IP
	if errClient != nil {
		return "", errClient
	}
	// 程序在使用完 response 后必须关闭 response 的主体。
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(responseClient.Body)

	body, _ := ioutil.ReadAll(responseClient.Body)
	clientIP := fmt.Sprintf("%s", string(body))
	return clientIP, nil
}

func Map2state(m map[string]interface{}) *amazonsChess.State {
	board, _ := toIntSlice(m["board"])
	current_player := int(m["current_player"].(float64))
	state := amazonsChess.NewState(&board, current_player)
	return state
}
