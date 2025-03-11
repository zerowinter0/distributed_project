//go:generate protoc --go_out=. gen/message.proto
package example_pkg

import (
	"bufio"
	"errors"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"google.golang.org/protobuf/proto"
)

//func readMsg(name string, reader *bufio.Reader, msg proto.Message) error {
//	//todo:延迟or未收到
//	// 读取消息长度
//	lenBuf := make([]byte, 4)
//	_, err := io.ReadFull(reader, lenBuf)
//	if err != nil {
//		if err != io.EOF {
//			log.Println(name, "error reading message length:", err.Error())
//		}
//		return err
//	}
//
//	msgLen := int(lenBuf[0])<<24 | int(lenBuf[1])<<16 | int(lenBuf[2])<<8 | int(lenBuf[3])
//	msgBuf := make([]byte, msgLen)
//
//	// 读取消息内容
//	_, err = io.ReadFull(reader, msgBuf)
//	if err != nil {
//		log.Println(name, "error reading message content:", err.Error())
//		return err
//	}
//
//	// 反序列化Protobuf消息
//
//	err = proto.Unmarshal(msgBuf, msg)
//	if err != nil {
//		log.Println(name, "error unmarshalling message:", err.Error())
//		return err
//	}
//	return nil
//}

func readMsg(name string, reader *bufio.Reader, msg proto.Message) error {
	// --- 读取配置文件 ---
	var delay int
	var noResponse bool
	var configPath = name + "_tst.cfg"

	// 打开配置文件
	file, err := os.Open(configPath)
	if err != nil {
		log.Printf("%s: 无法读取配置文件 %s，使用默认值 (delay=0, no_response=false)", name, configPath)
		delay = 0
		noResponse = false
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)

		// 读取延迟秒数（第一行）
		if scanner.Scan() {
			if val, err := strconv.Atoi(scanner.Text()); err == nil {
				delay = val
			} else {
				log.Printf("%s: 配置文件中延迟值解析失败，使用默认值 0", name)
			}
		}

		// 读取不应答标志（第二行）
		if scanner.Scan() {
			if val, err := strconv.Atoi(scanner.Text()); err == nil {
				noResponse = (val != 0)
			} else {
				log.Printf("%s: 配置文件中不应答标志解析失败，使用默认值 false", name)
			}
		}

		// 检查扫描错误
		if err := scanner.Err(); err != nil {
			log.Printf("%s: 配置文件扫描错误: %v", name, err)
		}
	}

	// --- 应用延迟 ---
	if noResponse {
		log.Printf("%s: 不应答模式，跳过消息处理", name)
		return errors.New("no response mode enabled") // 返回错误让调用者捕捉
	}

	time.Sleep(time.Duration(delay) * time.Second)

	// 读取消息长度
	lenBuf := make([]byte, 4)
	_, err = io.ReadFull(reader, lenBuf)
	if err != nil {
		//if err != io.EOF {
		//	log.Println(name, "error reading message length:", err.Error())
		//}
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
