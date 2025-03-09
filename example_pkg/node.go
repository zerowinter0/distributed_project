package example_pkg

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"log"
	"sync"
	"time"
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
	}
	runServer(name, thisAddr.Port)
}

func getInputAndSend(
	name string,
	name2addr map[string]Address,
) {
	message := fmt.Sprintf("node %s message", name)
	log.Println("node", name, "sending message:", message)

	wg := sync.WaitGroup{}
	time.Sleep(time.Duration(10) * time.Second)

	//发消息给所有其它节点
	for namePeer, ipAndPort := range name2addr {
		if namePeer != name {
			addr := fmt.Sprintf("%s:%d", ipAndPort.IP, ipAndPort.Port)
			wg.Add(1)
			go func() {
				defer wg.Done()
				runClient(name, namePeer, addr, message)
			}()

		}
	}
	wg.Wait()
}

func runClient(name string, namePeer string, address string, message string) {
	err := client(name, namePeer, address, message)
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
