package tcp

import (
	"encoding/binary"
	"net"
	"syscall"
	"testing"
	"unsafe"
)

//func TestHashAllinOne(t *testing.T) {
//	// 创建一个 TCP 连接
//	conn, err := net.Dial("tcp", "github.com:80")
//	if err != nil {
//		panic(err)
//	}
//	defer conn.Close()
//
//	// 获取 TCP 连接的本地地址和远程地址
//	tcpConn := conn.(*net.TCPConn)
//	localAddr := tcpConn.LocalAddr().(*net.TCPAddr)
//	remoteAddr := tcpConn.RemoteAddr().(*net.TCPAddr)
//
//	// 获取 TCP 连接的状态
//	state, err := getTCPState(tcpConn)
//	if err != nil {
//		panic(err)
//	}
//
//	// 将 TCP 连接状态添加到哈希表中
//	statMap := NewTCPStatMap()
//	statMap.Add(remoteAddr, localAddr, state)
//
//	// 输出当前 ESTABLISHED 状态的 TCP 连接数
//	count := statMap.Count(StateEstablished)
//	fmt.Printf("ESTABLISHED: %d\n", count)
//
//	// 关闭 TCP 连接
//	conn.Close()
//
//	// 从哈希表中删除 TCP 连接状态
//	statMap.Remove(remoteAddr, localAddr)
//}

func TestGetTCPState(t *testing.T) {
	// 创建一个本地监听的 TCP 服务器
	listener, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	// 启动一个协程，等待客户端连接
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
	}()

	// 在协程中连接服务器
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// 获取 TCP 连接的状态
	tcpConn := conn.(*net.TCPConn)
	state, err := getTCPState(tcpConn)
	if err != nil {
		t.Fatal(err)
	}

	// 断言获取到的状态是 StateEstablished
	if state != StateEstablished {
		t.Errorf("got state %v, want %v", state, StateEstablished)
	}
}

func TestGetSockOpt(t *testing.T) {
	// 创建一个本地监听的 TCP 服务器
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()

	// 启动一个协程，等待客户端连接
	go func() {
		conn, err := listener.Accept()
		if err != nil {
			t.Error(err)
			return
		}
		defer conn.Close()
	}()

	// 在协程中连接服务器
	conn, err := net.Dial("tcp", listener.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// 转换为 *net.TCPConn 类型
	tcpConn := conn.(*net.TCPConn)

	// 获取连接的文件描述符
	f, err := tcpConn.File()
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	// 调用 getSockOpt 函数
	var state uint32
	buf := make([]byte, 4)
	vallen := uintptr(4)
	err = getSockOpt(f.Fd(), syscall.IPPROTO_TCP, syscall.TCP_INFO, unsafe.Pointer(&buf[0]), &vallen)
	if err != nil {
		t.Fatalf("getSockOpt failed: %v", err)
	}

	state = binary.LittleEndian.Uint32(buf)
	t.Logf("TCP state: %v", state)
}
