package tcp

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"strconv"
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

	fd, err := unix.Socket(unix.AF_INET, unix.SOCK_STREAM, unix.IPPROTO_TCP)
	if err != nil {
		logrus.Warn("Failed to create socket: ", err)
		return err
	}

	ip, port, _ := net.SplitHostPort(cr.IPAddr)
	p, _ := strconv.Atoi(port)

	addr := unix.SockaddrInet4{Port: p}

	copy(addr.Addr[:], net.ParseIP(ip).To4())
	err = unix.Bind(fd, &addr)
	if err != nil {
		return err
	}

	err = unix.SetsockoptInt(fd, unix.IPPROTO_TCP, unix.TCP_DEFER_ACCEPT, 1)
	if err != nil {
		return err
	}

	err = unix.Listen(fd, unix.SOMAXCONN)
	if err != nil {
		return err
	}

	cr.fd = &fd
	return nil
}

// AcceptToFD 返回socket的文件描述符，后续零拷贝后无需再次进行转换
func (cr *CustomRule) AcceptToFD() (*os.File, error) {
	var cfd int
	var err error

	// 当 Accept 方法完成后，代表连接已经进入了传输层且已经连接建立完成
	cfd, _, err = unix.Accept(*cr.fd)
	if err != nil {
		return nil, fmt.Errorf("accept error: %v", err)
	}

	err = unix.SetNonblock(cfd, true) // 设置连接为非阻塞模式
	if err != nil {
		unix.Close(cfd)
		return nil, err
	}

	f := os.NewFile(uintptr(cfd), "")
	return f, nil
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
			if tcp.SYN && tcp.ACK {
				logrus.Warn(
					"malicious packets are detected and the source ip has been added to the block list (SYN+ACK): ",
					ip.SrcIP.String(),
				)
				gso.bloom.AddString(ip.SrcIP.String())
			}
		}
	}
}

// 达到了阈值，开始进行收敛
