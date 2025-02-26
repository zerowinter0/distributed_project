package example_pkg

import (
	"fmt"
	"testing"
)

func TestNode(t *testing.T) {
	n := 10 // N个节点
	name2addr := make(
		map[string]struct {
			ip   string
			port uint16
		})
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("n_%d", i+1)
		ip := "127.0.0.1"
		port := uint16(8080) + uint16(i) + 1
		name2addr[name] = struct {
			ip   string
			port uint16
		}{ip: ip, port: port}
	}

	chDone := make(chan bool, 1)
	for name, _ := range name2addr {
		//只有n1接收输入
		isClient := name == "n_1"
		go func(nodeName string,
			name2addr map[string]struct {
				ip   string
				port uint16
			},
			isClient bool) {
			node(name, name2addr, isClient, &chDone)
		}(name, name2addr, isClient)
	}
	_ = <-chDone
}
