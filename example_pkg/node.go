package example_pkg

import (
	"fmt"
	"log"
	"sync"
)

func node(
	name string,
	name2addr map[string]struct {
		ip   string
		port uint16
	},
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
	runServer(name, thisAddr.port)
}

func getInputAndSend(name string,
	name2addr map[string]struct {
		ip   string
		port uint16
	},
) {
	message := fmt.Sprintf("node %s message", name)
	log.Println("node", name, "sending message:", message)

	wg := sync.WaitGroup{}
	//发消息给所有其它节点
	for namePeer, ipAndPort := range name2addr {
		if namePeer != name {
			addr := fmt.Sprintf("%s:%d", ipAndPort.ip, ipAndPort.port)
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
