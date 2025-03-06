package example_pkg

import (
	"bufio"
	"example_pkg/gen"
	"fmt"
	"net"
	"time"
)

func client(name string, namePeer string, address string, message string, query_type int) error {
	var conn net.Conn
	var err error
	for {
		conn, err = net.Dial("tcp", address)
		if err != nil {
			time.Sleep(time.Duration(1000) * time.Millisecond)
		} else {
			break
		}
	}

	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	// 创建Protobuf消息
	msg := &gen.MyMessage{Content: message}

	watchAppendMessage(fmt.Sprintf("%s %s", name, msg))
	err = writeMsg(name, writer, msg)
	if err != nil {
		return err
	}

	// 读取响应
	err = readMsg(name, reader, msg)
	if err != nil {
		return err
	}
	if query_type == 1 {
		fmt.Println("node " + name + ": " + msg.Content)
	}
	return nil
}
