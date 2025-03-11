package example_pkg

import (
	"bufio"
	"example_pkg/gen"
	"fmt"
	"net"
	"time"
)

//	func client(name string, namePeer string, address string, message string, query_type int) error {
//		var conn net.Conn
//		var err error
//		for {
//			// 连接到服务器
//			conn, err = net.Dial("tcp", address)
//			// 如果连接失败，1000ms后重试，直到成功
//			if err != nil {
//				time.Sleep(time.Duration(1000) * time.Millisecond)
//			} else {
//				break
//			}
//		}
//
//		defer func(conn net.Conn) {
//			_ = conn.Close()
//		}(conn)
//
//		reader := bufio.NewReader(conn)
//		writer := bufio.NewWriter(conn)
//
//		// 创建Protobuf消息
//		msg := &gen.MyMessage{Content: message}
//
//		// 发送消息
//		watchAppendMessage(fmt.Sprintf("%s %s", name, msg))
//		// 将消息内容回显给客户端
//		err = writeMsg(name, writer, msg)
//		if err != nil {
//			return err
//		}
//
//		// 读取响应
//		err = readMsg(name, reader, msg)
//		if err != nil {
//			return err
//		}
//		// 处理响应
//		if query_type == 1 {
//			fmt.Println("node " + name + ": " + msg.Content)
//		}
//		return nil
//	}
const (
	maxRetries    = 300             // 最大重试次数
	retryInterval = 1 * time.Second // 重试间隔
	readTimeout   = 5 * time.Second // 读取超时时间
)

func client(name string, namePeer string, address string, message string, query_type int) error {
	var conn net.Conn
	var err error
	for retry := 0; retry < maxRetries; retry++ {
		// 连接到服务器
		conn, err = net.Dial("tcp", address)
		if err != nil {
			if retry < maxRetries-1 {
				time.Sleep(retryInterval)
				continue
			}
			return fmt.Errorf("连接失败，超过重试次数: %v", err)
		}

		// 设置读取超时
		conn.SetReadDeadline(time.Now().Add(readTimeout))

		reader := bufio.NewReader(conn)
		writer := bufio.NewWriter(conn)

		// 创建Protobuf消息
		msg := &gen.MyMessage{Content: message}

		// 发送消息
		watchAppendMessage(fmt.Sprintf("%s %s", name, msg))
		if err := writeMsg(name, writer, msg); err != nil {
			conn.Close()
			if retry < maxRetries-1 {
				time.Sleep(retryInterval)
				continue
			}
			return fmt.Errorf("发送失败，超过重试次数: %v", err)
		}

		// 读取响应
		if err := readMsg(name, reader, msg); err != nil {
			conn.Close()
			if retry < maxRetries-1 {
				time.Sleep(retryInterval)
				continue
			}
			return fmt.Errorf("读取响应失败，超过重试次数: %v", err)
		}

		// 成功处理响应
		conn.Close()
		if query_type == 1 {
			fmt.Println("node " + name + ": " + msg.Content)
		}
		return nil
	}
	return fmt.Errorf("未知错误")
}
