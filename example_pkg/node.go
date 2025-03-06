package example_pkg

import (
	"fmt"
	"log"

	"github.com/BurntSushi/toml"
)

func NewNode(cfgPath string) {
	var cfg Cfg
	fmt.Println("read cfg from path:", cfgPath)
	if _, err := toml.DecodeFile(cfgPath, &cfg); err != nil {
		panic(err)
	}

	node(cfg.Name, cfg.Name2Addr(), cfg.IsClient, nil)
}

func node(
	name string,
	name2addr map[string]Address,
	client bool,
	ch *chan bool,
) {
	thisAddr := name2addr[name]
	if client {
		go func() {
			getInputAndSend(name, name2addr)
			if *ch != nil {
				*ch <- true
			}
		}()
	} else {
		runServer(name, thisAddr.Port)
	}
}

func getInputAndSend(
	name string,
	name2addr map[string]Address,
) {
	master_map := NewSafeMap(name)

	fmt.Println("master node waiting for input...")

	serial_number := master_map.max_serial_number

	for {
		var key string
		var value string
		fmt.Scanln(&key, &value)
		var msg string
		query_type := 0
		if key == "query" {
			msg = encode_request(value)
			query_type = 1
			master_ans, _ := master_map.Query(value)
			fmt.Println("master node answer:" + master_ans)
		} else {
			msg = encode_kv(key, value, serial_number)
			master_map.Insert(key, value, serial_number)
		}

		serial_number += 1

		for namePeer, ipAndPort := range name2addr {
			if namePeer != name {
				addr := fmt.Sprintf("%s:%d", ipAndPort.IP, ipAndPort.Port)

				go func() {
					runClient(name, namePeer, addr, msg, query_type)
				}()

			}
		}
	}
}

func runClient(name string, namePeer string, address string, message string, query_type int) {
	err := client(name, namePeer, address, message, query_type)
	if err != nil {
		log.Fatal(err)
	}
}

func runServer(name string, port uint16) {
	err := server(name, port)
	if err != nil {
		log.Fatal(err)
	}
}
