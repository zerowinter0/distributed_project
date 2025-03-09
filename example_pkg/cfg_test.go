//package example_pkg
//
//import (
//	"github.com/BurntSushi/toml"
//	"os"
//	"testing"
//)
//
//func TestCfg(t *testing.T) {
//	cfg := Cfg{
//		Name: "node1",
//		Peers: []Peer{
//			{Name: "node1", Ip: "192.168.0.1", Port: 8081},
//			{Name: "node2", Ip: "192.168.0.2", Port: 8082},
//			{Name: "node3", Ip: "192.168.0.3", Port: 8083},
//		},
//		IsClient: true,
//	}
//	file, err := os.OpenFile("/tmp/cfg.toml", os.O_WRONLY|os.O_CREATE, 0600)
//	if err != nil {
//		panic(err)
//	}
//	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
//		panic(err)
//	}
//	err = file.Close()
//	if err != nil {
//		panic(err)
//	}
//	if _, err := toml.DecodeFile("/tmp/cfg.toml", &cfg); err != nil {
//		panic(err)
//	}
//}

package example_pkg

import (
	"github.com/BurntSushi/toml"
	"os"
	"path/filepath"
	"testing"
)

func TestCfg(t *testing.T) {
	// 定义配置结构体
	cfg := Cfg{
		Name: "node1",
		Peers: []Peer{
			{Name: "node1", Ip: "192.168.0.1", Port: 8081},
			{Name: "node2", Ip: "192.168.0.2", Port: 8082},
			{Name: "node3", Ip: "192.168.0.3", Port: 8083},
		},
		IsClient: true,
	}

	// 获取系统临时目录路径，并拼接文件名
	tempDir := os.TempDir()
	filePath := filepath.Join(tempDir, "cfg.toml")

	// 创建并写入配置文件
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		t.Fatalf("failed to create config file: %v", err)
	}
	defer func() {
		// 测试完成后删除临时文件
		if err := os.Remove(filePath); err != nil {
			t.Logf("failed to remove temp file: %v", err)
		}
	}()

	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		t.Fatalf("failed to encode config to file: %v", err)
	}
	if err := file.Close(); err != nil {
		t.Fatalf("failed to close file: %v", err)
	}

	// 从文件中解码配置
	var decodedCfg Cfg
	if _, err := toml.DecodeFile(filePath, &decodedCfg); err != nil {
		t.Fatalf("failed to decode config from file: %v", err)
	}

	// 验证解码后的配置是否正确
	if decodedCfg.Name != cfg.Name || len(decodedCfg.Peers) != len(cfg.Peers) || decodedCfg.IsClient != cfg.IsClient {
		t.Errorf("decoded config does not match original config")
	}
}
