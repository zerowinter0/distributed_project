package example_pkg

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestWatchServer(t *testing.T) {
	server := NewWatchServer(19099)
	ch := make(chan string)
	m := make(map[string]int)
	server.Register("/trace", &ch, &m,
		func(writer http.ResponseWriter, request *http.Request) {
			uri := request.RequestURI
			s := server.context[uri]
			if s != nil {
				maps := reflect.ValueOf(s).Interface().(*map[string]int)
				if maps != nil {
					for k, v := range *maps {
						_, err := writer.Write([]byte(fmt.Sprintf("message:%s, count:%d\n", k, v)))
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
				maps := reflect.ValueOf(object).Interface().(*map[string]int)
				if maps != nil && channel != nil {
					for {
						s := <-*channel
						if n, ok := (*maps)[s]; ok {
							(*maps)[s] = n + 1
						} else {
							(*maps)[s] = 1
						}
					}
				}
			}
		},
	)
	go server.Serve()
	func(ch *chan string) {
		for j := 0; j < 2; j++ {
			for i := 0; i < 2; i++ {
				time.Sleep(1 * time.Second)
				*ch <- fmt.Sprintf("  hello %d", i)
			}
		}
	}(&ch)

}
