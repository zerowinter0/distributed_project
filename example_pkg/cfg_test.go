package example_pkg

import (
	"github.com/BurntSushi/toml"
	"os"
	"testing"
)

func TestCfg(t *testing.T) {
	cfg := Cfg{
		Name: "node1",
		Peers: []Peer{
			{Name: "node1", Ip: "192.168.0.1", Port: 8081},
			{Name: "node2", Ip: "192.168.0.2", Port: 8082},
			{Name: "node3", Ip: "192.168.0.3", Port: 8083},
		},
		IsClient: true,
	}
	file, err := os.OpenFile("/tmp/cfg.toml", os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		panic(err)
	}
	if err := toml.NewEncoder(file).Encode(cfg); err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
	if _, err := toml.DecodeFile("/tmp/cfg.toml", &cfg); err != nil {
		panic(err)
	}
}
