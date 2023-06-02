package tcp

// 使用带计数器的哈希表来缓存全部TCP连接状态
// 以便快速判断是否需要拦截或是进行收敛

import (
	"encoding/binary"
	"fmt"
	"github.com/spaolacci/murmur3"
	"net"
	"sync"
	"syscall"
	"unsafe"
)

// State 表示 TCP 连接状态
type State int

const (
	// StateUnknown 表示未知状态
	StateUnknown State = iota

	// StateEstablished 表示已建立连接
	StateEstablished

	// StateSynSent 表示 SYN 已发送
	StateSynSent

	// StateSynReceived 表示 SYN 已接收
	StateSynReceived

	// StateFinWait1 表示等待远程发送 FIN
	StateFinWait1

	// StateFinWait2 表示等待远程发送 FIN，同时自己也发送了 FIN
	StateFinWait2

	// StateTimeWait 表示等待定时器超时
	StateTimeWait

	// StateCloseWait 表示等待远程关闭连接
	StateCloseWait

	// StateLastAck 表示等待远程发送 ACK
	StateLastAck

	// StateClosing 表示同时发送 FIN，等待远程响应
	StateClosing

	// StateClosed 表示连接已关闭
	StateClosed
)

// StateString 返回 StateEstablished 连接状态的字符串表示
func StateString(state State) string {
	switch state {
	case StateEstablished:
		return "ESTABLISHED"
	case StateSynSent:
		return "SYN_SENT"
	case StateSynReceived:
		return "SYN_RECEIVED"
	case StateFinWait1:
		return "FIN_WAIT1"
	case StateFinWait2:
		return "FIN_WAIT2"
	case StateTimeWait:
		return "TIME_WAIT"
	case StateCloseWait:
		return "CLOSE_WAIT"
	case StateLastAck:
		return "LAST_ACK"
	case StateClosing:
		return "CLOSING"
	case StateClosed:
		return "CLOSED"
	default:
		return "UNKNOWN"
	}
}

// Stat 表示 TCP 连接状态和计数器
type Stat struct {
	State State // TCP 连接状态
	Count int   // 计数器
}

// TCPStatMap 表示带有状态计数器的 TCP 连接状态哈希表
type TCPStatMap struct {
	sync.RWMutex
	m map[uint64]Stat // 哈希表
}

// NewTCPStatMap 返回一个新的 TCPStatMap 对象
func NewTCPStatMap() *TCPStatMap {
	return &TCPStatMap{
		m: make(map[uint64]Stat),
	}
}

// Add 添加一个 TCP 连接状态
func (m *TCPStatMap) Add(remoteAddr, localAddr *net.TCPAddr, state State) {
	// 将 remoteAddr 和 localAddr 拼接成一个字符串，并计算哈希值
	key := hashRemoteLocalAddr(remoteAddr, localAddr)

	// 获取写锁
	m.Lock()
	defer m.Unlock()

	// 查找哈希表中是否存在对应的 key
	stat, ok := m.m[key]
	if !ok {
		// 如果不存在，则新增一个 key
		stat = Stat{State: state, Count: 1}
	} else {
		stat.State = state
		stat.Count++
	}

	// 更新哈希表
	m.m[key] = stat
}

// Remove 删除一个 TCP 连接状态
func (m *TCPStatMap) Remove(remoteAddr, localAddr *net.TCPAddr) {
	// 将 remoteAddr 和 localAddr 拼接成一个字符串，并计算哈希值
	key := hashRemoteLocalAddr(remoteAddr, localAddr)

	// 获取写锁
	m.Lock()
	defer m.Unlock()

	// 从哈希表中删除对应的 key
	delete(m.m, key)
}

// Count 返回指定状态的连接数
func (m *TCPStatMap) Count(state State) int {
	// 获取读锁
	m.RLock()
	defer m.RUnlock()

	// 遍历哈希表，统计指定状态的连接数
	count := 0
	for _, stat := range m.m {
		if stat.State == state {
			count += stat.Count
		}
	}
	return count
}

// hashRemoteLocalAddr 将 remoteAddr 和 localAddr 拼接成一个字符串，并计算哈希值
func hashRemoteLocalAddr(remoteAddr, localAddr *net.TCPAddr) uint64 {
	s := fmt.Sprintf("%s:%d-%s:%d", remoteAddr.IP.String(), remoteAddr.Port, localAddr.IP.String(), localAddr.Port)
	h := murmur3.New64()
	h.Write([]byte(s))
	return h.Sum64()
}

func getSockOpt(fd uintptr, level, name int, val unsafe.Pointer, vallen *uintptr) error {
	_, _, errno := syscall.Syscall6(
		syscall.SYS_GETSOCKOPT,
		fd,
		uintptr(level),
		uintptr(name),
		uintptr(val),
		uintptr(unsafe.Pointer(vallen)),
		0,
	)

	if errno != 0 {
		return error(errno)
	}
	return nil
}

// getTCPState 获取 TCP 连接的状态
// 获取 TCP 连接的状态
func getTCPState(conn *net.TCPConn) (State, error) {
	f, err := conn.File()
	if err != nil {
		return StateUnknown, err
	}
	// 获取 TCP 连接的本地地址和远程地址
	localAddr := conn.LocalAddr().(*net.TCPAddr)
	remoteAddr := conn.RemoteAddr().(*net.TCPAddr)

	// 调用 getSockOpt 获取 TCP 连接状态
	var state uint32
	buf := make([]byte, 4)
	vallen := uintptr(4) // 创建一个 uintptr 类型的变量来存储 vallen 的值
	err = getSockOpt(f.Fd(), syscall.IPPROTO_TCP, syscall.TCP_INFO, unsafe.Pointer(&buf[0]), &vallen)
	if err != nil {
		return StateUnknown, err
	}
	state = binary.LittleEndian.Uint32(buf)

	// 将获取到的状态添加到 TCPStatMap
	statMap := NewTCPStatMap()
	statMap.Add(remoteAddr, localAddr, State(state))
	defer statMap.Remove(remoteAddr, localAddr)

	return State(state), nil
}
