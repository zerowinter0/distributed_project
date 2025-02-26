//go:generate protoc --go_out=. gen/message.proto
package example_pkg

import (
	"bufio"
	"google.golang.org/protobuf/proto"
	"io"
	"log"
)

func readMsg(name string, reader *bufio.Reader, msg proto.Message) error {
	// 读取消息长度
	lenBuf := make([]byte, 4)
	_, err := io.ReadFull(reader, lenBuf)
	if err != nil {
		if err != io.EOF {
			log.Println(name, "error reading message length:", err.Error())
		}
		return err
	}

	msgLen := int(lenBuf[0])<<24 | int(lenBuf[1])<<16 | int(lenBuf[2])<<8 | int(lenBuf[3])
	msgBuf := make([]byte, msgLen)

	// 读取消息内容
	_, err = io.ReadFull(reader, msgBuf)
	if err != nil {
		log.Println(name, "error reading message content:", err.Error())
		return err
	}

	// 反序列化Protobuf消息

	err = proto.Unmarshal(msgBuf, msg)
	if err != nil {
		log.Println(name, "error unmarshalling message:", err.Error())
		return err
	}
	return nil
}

func writeMsg(name string, writer *bufio.Writer, msg proto.Message) error {
	// 将消息内容回显给客户端

	buf, err := proto.Marshal(msg)
	if err != nil {
		log.Println(name, "error marshalling response:", err.Error())
		return err
	}

	// 发送消息长度
	lenBuf := make([]byte, 4)
	msgLen := len(buf)
	lenBuf[0] = byte(msgLen >> 24)
	lenBuf[1] = byte(msgLen >> 16)
	lenBuf[2] = byte(msgLen >> 8)
	lenBuf[3] = byte(msgLen)
	_, err = writer.Write(lenBuf)
	if err != nil {
		log.Println(name, "error writing response length:", err.Error())
		return err
	}

	// 发送消息内容
	_, err = writer.Write(buf)
	if err != nil {
		log.Println(name, "error writing response content:", err.Error())
		return err
	}

	err = writer.Flush()
	if err != nil {
		return err
	}
	return nil
}
