package example_pkg

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func watchServer(port uint16) WatchContext {
	server := NewWatchServer(port)
	context := WatchContext{
		channel:  make(chan string),
		messages: make([]string, 0),
	}
	server.Register("/trace_message", &context.channel, &context,
		func(writer http.ResponseWriter, request *http.Request) {
			uri := request.RequestURI
			s := server.context[uri]
			if s != nil {
				ctx := reflect.ValueOf(s).Interface().(*WatchContext)
				if ctx != nil {
					for _, m := range ctx.messages {
						_, err := writer.Write([]byte(fmt.Sprintf("%s\n", m)))
						if err != nil {
							log.Fatal(err)
						}
					}
				}
			}
		},
		func(channel interface{}, object interface{}) {
			if channel != nil && object != nil {
				channel := reflect.ValueOf(channel).Interface().(*chan string)
				ctx := reflect.ValueOf(object).Interface().(*WatchContext)
				if ctx != nil && channel != nil {
					for {
						s := <-*channel
						ctx.messages = append(ctx.messages, s)
					}
				}
			}
		},
	)
	go server.Serve()
	return context
}

func TestNode(t *testing.T) {
	n := 10 // N个节点
	context := watchServer(9099)
	setWatchCtx(&context)

	name2addr := make(
		map[string]Address)
	for i := 0; i < n; i++ {
		name := fmt.Sprintf("n_%d", i+1)
		ip := "127.0.0.1"
		port := uint16(12000) + uint16(i) + 1
		name2addr[name] = Address{
			IP:   ip,
			Port: port,
		}
	}

	chDone := make(chan bool, 1)
	for name := range name2addr {
		//只有n1接收输入
		if name == "n_1" {
			continue
		}
		go func(nodeName string,
			name2addr map[string]Address,
			isClient bool) {
			node(name, name2addr, isClient, &chDone)
		}(name, name2addr, false)
	}
	for name := range name2addr {
		if name == "n_1" {
			node(name, name2addr, true, &chDone)
		}
	}
	_ = <-chDone
	time.Sleep(time.Duration(60) * time.Second)
}
