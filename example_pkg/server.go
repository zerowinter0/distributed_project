package example_pkg

import (
	"bufio"
	"example_pkg/gen"
	"fmt"
	"io"
	"log"
	"net"
)

func server(name string, port uint16) error {

	kv_map := NewSafeMap(name)

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println("Error listening:", err.Error())
		return err
	}
	defer func(listen net.Listener) {
		_ = listen.Close()
	}(listen)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("Error accepting:", err.Error())
			continue
		}

		go handleClient(name, conn, kv_map)
	}
}

func handleClient(name string, conn net.Conn, kv_map *SafeMap) {
	err := _handleClient(name, conn, kv_map)
	if err != nil {
		log.Println(name, "error handling client:", err.Error())
	}
}

func _handleClient(name string, conn net.Conn, kv_map *SafeMap) error {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		msg := &gen.MyMessage{}
		err := readMsg(name, reader, msg)
		if err != nil {
			if err.Error() == "no response mode enabled" {
				// 处理不应答模式的情况
				log.Println(name, "检测到不应答模式，跳过处理")
				return nil // 或其他逻辑
			}
			if err == io.EOF {
				return nil
			}
			return err
		}

		log.Println(name, "received message: ", msg.Content)
		request_type, kv_data, err := decode_request(msg.Content)

		if err != nil {
			return err
		}
		if request_type == 0 {
			kv_map.Insert(kv_data.key, kv_data.value.value, kv_data.value.serial_number)
		} else if request_type == 1 {
			value, exist := kv_map.Query(kv_data.key)
			if exist {
				msg.Content = value
			} else {
				msg.Content = ""
			}
		}

		watchAppendMessage(fmt.Sprintf("%s %s", name, msg))
		// 将消息内容回显给客户端
		response := &gen.MyMessage{Content: msg.Content}
		err = writeMsg(name, writer, response)
		if err != nil {
			return err
		}
	}
}
