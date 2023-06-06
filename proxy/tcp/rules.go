package tcp

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
	"net"
	"os"
	"strconv"
	"syscall"
)

// CustomRule 重载连接行为，用于实现底层的自定义连接规则
type CustomRule struct {
	net.Listener
	IPAddr string
	fd     *int
}

func (cr *CustomRule) Listen() error {
	if len(cr.IPAddr) == 0 {
		logrus.Errorf("ipaddr must have value")
	}

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		logrus.Warn("Failed to create socket: ", err)
		return err
	}

	ip, port, _ := net.SplitHostPort(cr.IPAddr)
	p, _ := strconv.Atoi(port)

	addr := syscall.SockaddrInet4{Port: p}

	copy(addr.Addr[:], net.ParseIP(ip).To4())
	err = syscall.Bind(fd, &addr)
	if err != nil {
		return err
	}

	err = syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_DEFER_ACCEPT, 1)
	if err != nil {
		return err
	}

	err = syscall.Listen(fd, syscall.SOMAXCONN)
	if err != nil {
		return err
	}

	cr.fd = &fd
	return nil
}

func (cr *CustomRule) Accept() (net.Conn, error) {
	var cfd int
	var err error

	// 当 Accept 方法完成后，代表连接已经进入了传输层且已经连接建立完成
	cfd, _, err = syscall.Accept(*cr.fd)
	if err != nil {
		return nil, fmt.Errorf("accept error: %v", err)
	}

	err = syscall.SetNonblock(cfd, true) // 设置连接为非阻塞模式
	if err != nil {
		syscall.Close(cfd)
		return nil, err
	}

	f := os.NewFile(uintptr(cfd), "")
	defer f.Close()

	conn, err := net.FileConn(f)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func LoadDenyIPRules(gso *ServerOptions, ips []string) {
	for _, v := range ips {
		gso.bloom.AddString(v)
	}
}

// 对于基于ACK计数器的防御

// 对于大量SYN_RECV状态的缓解

// WithTCPServerSYNACKRecv 对于服务端接受到异常握手包的拦截
func WithTCPServerSYNACKRecv(p <-chan gopacket.Packet, gso *ServerOptions) {
	for {
		select {
		case packet, ok := <-p:
			if !ok {
				logrus.Warn("packet channel closed")
				return
			}
			ipLayer := packet.Layer(layers.LayerTypeIPv4)
			if ipLayer == nil {
				continue
			}

			ip, _ := ipLayer.(*layers.IPv4)
			tcpLayer := packet.Layer(layers.LayerTypeTCP)
			if tcpLayer == nil {
				continue
			}

			tcp := tcpLayer.(*layers.TCP)
			logrus.Println(tcp.SYN, tcp.ACK)
			if tcp.SYN && tcp.ACK {
				logrus.Warn("SYN ACK Recv: ", ip.SrcIP.String())
				gso.bloom.AddString(ip.SrcIP.String())
			}
		}
	}
}

// 达到了阈值，开始进行收敛
