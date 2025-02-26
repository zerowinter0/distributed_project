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
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Println("Error listening:", err.Error())
		return err
	}
	defer func(listen net.Listener) {
		_ = listen.Close()
	}(listen)

	log.Printf("Server %s is listening on port %d...\n", name, port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Println("Error accepting:", err.Error())
			continue
		}

		go handleClient(name, conn)
	}
}

func handleClient(name string, conn net.Conn) {
	err := _handleClient(name, conn)
	if err != nil {
		log.Println(name, "error handling client:", err.Error())
	}
}

func _handleClient(name string, conn net.Conn) error {
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	for {
		msg := &gen.MyMessage{}
		err := readMsg(name, reader, msg)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		log.Println(name, "received message: ", msg.Content)

		// 将消息内容回显给客户端
		response := &gen.MyMessage{Content: msg.Content}
		err = writeMsg(name, writer, response)
		if err != nil {
			return err
		}
	}
}
