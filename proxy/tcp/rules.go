package tcp

import (
	"github.com/sirupsen/logrus"
	"net"
	"strconv"
	"syscall"
)

// CustomRule 重载连接行为，用于实现底层的自定义连接规则
type CustomRule struct {
	net.Listener
	IPAddr net.Addr
}

func (cr *CustomRule) Accept() (net.Conn, error) {
	fd,err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		logrus.Warn("Failed to create socket: ", err)
		return nil, err
	}
	ip,port,_ := net.SplitHostPort(cr.Addr().String())
	p,_ := strconv.Atoi(port)

	addr := syscall.SockaddrInet4{Port: p}
	copy(addr.Addr[:], net.ParseIP(ip).To4())
	err = syscall.Bind(fd, &addr)
	if err != nil {
		return nil,err
	}

	err = syscall.Listen(fd, syscall.SOMAXCONN)
	if err != nil {
		return nil,err
	}

	return cr.Listener.Accept()
}

func (cr *CustomRule) checkConnHandshake(conn net.Conn) error {
	return nil
}

// 对于基于ACK计数器的防御

// 对于大量SYN_RECV状态的缓解

// 达到了阈值，开始进行收敛
