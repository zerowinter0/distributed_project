package example_pkg

type Cfg struct {
	Name     string `toml:"name"`
	Peers    []Peer `toml:"peers"`
	IsClient bool   `toml:"is_client"`
}

type Peer struct {
	Name string `toml:"name"`
	Ip   string `toml:"ip"`
	Port uint16 `toml:"port"`
}

func (c Cfg) Name2Addr() map[string]Address {
	maps := make(map[string]Address)
	for _, peer := range c.Peers {
		maps[peer.Name] = Address{IP: peer.Ip, Port: peer.Port}
	}
	return maps
}
