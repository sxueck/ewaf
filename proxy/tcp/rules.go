package tcp

import (
	"errors"
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
}

func (cr *CustomRule) Accept() (net.Conn, error) {
	if len(cr.IPAddr) == 0 {
		logrus.Errorf("ipaddr must have value")
	}

	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	if err != nil {
		logrus.Warn("Failed to create socket: ", err)
		return nil, err
	}

	ip, port, _ := net.SplitHostPort(cr.IPAddr)
	p, _ := strconv.Atoi(port)

	addr := syscall.SockaddrInet4{Port: p}

	logrus.Println(ip, port, 1)
	copy(addr.Addr[:], net.ParseIP(ip).To4())
	err = syscall.Bind(fd, &addr)
	if err != nil {
		return nil, err
	}

	err = syscall.SetsockoptInt(fd, syscall.IPPROTO_TCP, syscall.TCP_DEFER_ACCEPT, 1)
	if err != nil {
		return nil, err
	}

	err = syscall.Listen(fd, syscall.SOMAXCONN)
	if err != nil {
		return nil, err
	}

	var cfd int
	cfd, _, err = syscall.Accept(fd)
	if err != nil {
		panic(err)
	}

	err = syscall.SetNonblock(cfd, true) // 设置连接为非阻塞模式
	if err != nil {
		panic(err)
	}

	err = cr.checkConnHandshake(cfd) // 检测三次握手标识
	if err != nil {
		syscall.Close(cfd)
	}

	connFp := os.NewFile(uintptr(cfd), "")
	defer connFp.Close()

	conn, err := net.FileConn(connFp)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func (cr *CustomRule) checkConnHandshake(connFd int) error {
	bs := make([]byte, 1024)
	n, err := syscall.Read(connFd, bs)
	if err != nil {
		return err
	}

	if n < 3 || bs[0] != 0x16 || bs[1] != 0x03 || bs[2] != 0x01 {
		return errors.New("invalid handshake")
	}
	return nil
}

// 对于基于ACK计数器的防御

// 对于大量SYN_RECV状态的缓解

// 达到了阈值，开始进行收敛
