package main

import (
	"example_pkg"
	"fmt"
	"os"
	"strconv"
	"time"
)

func start_client_node(idx int) {
	path := fmt.Sprintf("./cfg_n%d.toml", idx)
	go example_pkg.NewNode(path)
}

func main() {
	args := os.Args

	// 打印程序名
	fmt.Println("程序名:", args[0])

	// 检查是否有额外参数
	if len(args) > 1 {
		for _, arg := range args[1:] {
			num, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Println("请输入一个整数,代表节点编号。1号为主节点,其他为从节点,且从节点每个都有自己的cfg.toml")
			} else {
				start_client_node(num)
			}
		}
	} else {
		start_client_node(1)
	}

	time.Sleep(100 * time.Second)
}
