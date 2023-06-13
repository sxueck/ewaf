package pkg

// generator-url: https://mholt.github.io/json-to-go/

type GlobalConfig struct {
	Global  Global   `json:"global"`
	Servers []Server `json:"servers"`
}

type Global struct {
	RstConn        int    `json:"rst_conn" mapstructure:"rst_conn"`
	Reuseport      bool   `json:"reuseport"`
	TunChannelCidr string `json:"tun_channel_cidr" mapstructure:"tun_channel_cidr"`
	TCPManagement  string `json:"tcp_management" mapstructure:"tcp_management"`
	Interface      string `json:"interface"` // 网口地址
}

type Frontend struct {
	Type       string     `json:"type"`
	HostName   string     `json:"host_name" mapstructure:"host_name"`
	ListenPort int        `json:"listen_port" mapstructure:"listen_port"`
	Location   []Location `json:"location"`
}

type Server struct {
	Frontend Frontend `json:"frontend"`
}

type Backend struct {
	ByPass string `json:"by_pass" mapstructure:"by_pass"`
}

type Location struct {
	Backend Backend `json:"backend"`
}
