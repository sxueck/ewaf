package tcp

// 使用带计数器的哈希表来缓存全部TCP连接状态
// 以便快速判断是否需要拦截或是进行收敛

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/spaolacci/murmur3"
	"log"
	"net"
	"sync"
	"syscall"
	"time"
	"unsafe"
)

// State 表示 TCP 连接状态
type State int

const NextStatusInterval = 2

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

type Stat struct {
	State State // TCP 连接状态
	Count int   // 计数器
}

type StatMap struct {
	sync.RWMutex
	m map[uint64]Stat // 哈希表
}

func NewTCPStatMap() *StatMap {
	return &StatMap{
		m: make(map[uint64]Stat),
	}
}

// Add 添加一个 TCP 连接状态
// 在哈希表中，TCP连接状态同时只能有一条，注意这里使用了五元组进行单一连接的判定
func (m *StatMap) Add(remoteAddr, localAddr *net.TCPAddr, state State) {
	key := hashRemoteLocalAddr(remoteAddr, localAddr)

	m.Lock()
	defer m.Unlock()

	stat, ok := m.m[key]
	if !ok {
		stat = Stat{State: state, Count: 1}
	} else {
		stat.State = state
		stat.Count++
	}

	// 更新哈希表
	m.m[key] = stat
}

// Remove 删除一个 TCP 连接状态
func (m *StatMap) Remove(remoteAddr, localAddr *net.TCPAddr) {
	key := hashRemoteLocalAddr(remoteAddr, localAddr)

	m.Lock()
	defer m.Unlock()

	delete(m.m, key)
}

// Count 返回指定状态的连接数
func (m *StatMap) Count(state State) int {
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

func hashRemoteLocalAddr(remoteAddr, localAddr *net.TCPAddr) uint64 {
	s := fmt.Sprintf("%s:%d-%s:%d",
		remoteAddr.IP.String(), remoteAddr.Port, localAddr.IP.String(), localAddr.Port)
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

// GetTCPState 获取 TCP 连接的状态
// 获取 TCP 连接的状态
func GetTCPState(conn *net.TCPConn, statMap *StatMap) error {
	f, err := conn.File()
	if err != nil {
		return err
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
		return err
	}
	state = binary.LittleEndian.Uint32(buf)

	statMap.Add(remoteAddr, localAddr, State(state))

	statMap.RWMutex.Lock()
	n := statMap.m[hashRemoteLocalAddr(remoteAddr, localAddr)].State
	statMap.RWMutex.Unlock()

	if n == StateUnknown || n == StateClosed || n == StateClosing {
		statMap.Remove(remoteAddr, localAddr)
	}

	return nil
}

func ContinuousGetTCPState(ctx context.Context, conn *net.Conn, statMap *StatMap) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.NewTicker(time.Second * time.Duration(NextStatusInterval)).C:
			err := GetTCPState((*conn).(*net.TCPConn), statMap)
			if err != nil {
				log.Println(err)
				return
			}
		}
	}
}
