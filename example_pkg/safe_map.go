package example_pkg

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Value struct {
	value         string
	serial_number int
}

type KV_pair struct {
	key   string
	value Value
}

func splitString(input string) (string, string, string) {
	parts := strings.Split(input, " ") // 按空格拆分字符串
	if len(parts) != 3 {
		return "", "", "" // 如果不是恰好两个部分，返回空字符串
	}
	return parts[0], parts[1], parts[2] // 返回两个子串
}

func encode_request(key string) string {
	// 定义一个 map
	data := map[string]interface{}{
		"type": 1, //0=insertion, 1=request
		"key":  key,
	}

	// 将 map 编码为 JSON 字符串
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("编码失败:", err)
	}
	return string(jsonData)
}

func encode_kv(key string, value string, serial_number int) string {
	// 定义一个 map
	data := map[string]interface{}{
		"type":          0, //0=insertion, 1=request
		"key":           key,
		"value":         value,
		"serial_number": serial_number,
	}

	// 将 map 编码为 JSON 字符串
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("编码失败:", err)
	}
	return string(jsonData)
}

func decode_request(data string) (int, KV_pair, error) {
	// 定义一个 map 来存储解码后的数据
	var kv_data map[string]interface{}

	// 将 JSON 字符串解码为 map
	err := json.Unmarshal([]byte(data), &kv_data)
	if err != nil {
		return 0, KV_pair{}, fmt.Errorf("解码失败: %v", err)
	}

	// 检查并提取 type
	TypeFloat, ok := kv_data["type"].(float64)
	if !ok {
		return 0, KV_pair{}, errors.New("type 不是数字类型")
	}
	request_type := int(TypeFloat)

	// 提取 key
	key, ok := kv_data["key"].(string)
	if !ok {
		return request_type, KV_pair{}, errors.New("key 不是字符串类型")
	}

	if request_type == 1 { //request key
		return request_type, KV_pair{
			key: key,
			value: Value{
				value:         key,
				serial_number: 0,
			},
		}, nil
	}

	// 提取 value
	value, ok := kv_data["value"].(string)
	if !ok {
		return request_type, KV_pair{}, errors.New("value 不是字符串类型")
	}

	// 提取 serial_number 并转换为 int
	serialNumberFloat, ok := kv_data["serial_number"].(float64)
	if !ok {
		return request_type, KV_pair{}, errors.New("serial_number 不是数字类型")
	}
	serialNumber := int(serialNumberFloat)

	// 构造并返回 KV_pair
	return request_type, KV_pair{
		key: key,
		value: Value{
			value:         value,
			serial_number: serialNumber,
		},
	}, nil
}

// 定义结构体
type SafeMap struct {
	data              map[string]Value // 内部存储数据的 map
	mu                sync.Mutex       // 互斥锁，用于保护 map 的并发访问
	filename          string
	max_serial_number int //用于master节点
}

// 新建 SafeMap 实例
func NewSafeMap(name string) *SafeMap {
	mp := &SafeMap{
		data:              make(map[string]Value),
		filename:          name + ".log",
		max_serial_number: 0,
	}

	filename := name + ".log"
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return nil
	}
	defer file.Close() // 确保文件关闭

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		s := scanner.Text()
		key, value, str_num := splitString(s)
		if key == "" && value == "" && str_num == "" {
			continue
		}
		num, _ := strconv.Atoi(str_num)
		mp.Recover(key, value, num)
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Reading error:", err)
		return nil
	}

	return mp
}

// 插入键值对
func (sm *SafeMap) Insert(key string, value string, serial_number int) {
	sm.mu.Lock()         // 加锁
	defer sm.mu.Unlock() // 确保在函数结束时解锁

	old_value, exists := sm.data[key]
	if exists && old_value.serial_number >= serial_number {
		//数据过时
		return
	}
	sm.data[key] = Value{value, serial_number}

	//写入文件
	file, err := os.OpenFile(sm.filename, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	defer file.Close()
	content := key + " " + value + " " + strconv.Itoa(serial_number)
	if _, err := file.WriteString("\n" + content); err != nil {
		fmt.Println("Appending error:", err)
		return
	}

}

// 用于重启恢复
func (sm *SafeMap) Recover(key string, value string, serial_number int) {
	sm.mu.Lock()         // 加锁
	defer sm.mu.Unlock() // 确保在函数结束时解锁

	old_value, exists := sm.data[key]
	if exists && old_value.serial_number >= serial_number {
		//数据过时
		return
	}
	sm.data[key] = Value{value, serial_number}
	if serial_number > sm.max_serial_number {
		sm.max_serial_number = serial_number
	}
}

// 查询键对应的值
func (sm *SafeMap) Query(key string) (string, bool) {
	sm.mu.Lock()         // 加锁
	defer sm.mu.Unlock() // 确保在函数结束时解锁

	value, exists := sm.data[key]
	if !exists {
		return "", exists
	}
	return value.value, exists
}
