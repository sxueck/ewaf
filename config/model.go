package config

// generator-url: https://mholt.github.io/json-to-go/

type cfg struct {
	Global  Global    `json:"global"`
	Servers []Servers `json:"servers"`
}

type Global struct {
	RstConn        int    `json:"rst_conn"`
	Reuseport      bool   `json:"reuseport"`
	TunChannelCidr string `json:"tun_channel_cidr"`
}

type Backend struct {
	ByPass string `json:"by_pass"`
}

type Location struct {
	Backend Backend `json:"backend"`
}

type Server struct {
	Type       string     `json:"type"`
	HostName   string     `json:"host_name"`
	ListenPort int        `json:"listen_port"`
	Location   []Location `json:"location"`
}

type Servers struct {
	Server Server `json:"server"`
}
