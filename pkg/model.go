package pkg

// generator-url: https://mholt.github.io/json-to-go/

type GlobalConfig struct {
	Global  Global    `json:"global"`
	Servers *[]Server `json:"servers"`
}

type Global struct {
	RstConn        int    `json:"rst_conn"`
	Reuseport      bool   `json:"reuseport"`
	TunChannelCidr string `json:"tun_channel_cidr"`
	TcpManagement  string `json:"tcp_management"`
}

type Frontend struct {
	Type       string      `json:"type"`
	HostName   string      `json:"host_name"`
	ListenPort int         `json:"listen_port"`
	Location   *[]Location `json:"location"`
}

type Server struct {
	Frontend Frontend `json:"frontend"`
}

type Backend struct {
	Method string `json:"method"`
	ByPass string `json:"by_pass"`
}

type Location struct {
	Backend Backend `json:"backend"`
}
