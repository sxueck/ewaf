//go:build linux

package water

import (
	"github.com/sirupsen/logrus"
	"github.com/vishvananda/netlink"
	"golang.org/x/sys/unix"
	"net"
	"os"
	"strings"
	"sync"
	"unsafe"
)

const (
	defaultMTU     = 1420
	ifReqSize      = unix.IFNAMSIZ + 0x40
	tunDevicePoint = "/dev/net/tun"
	TunLabel       = "ewaf-tun"
)

type TUNConfigure struct {
	tunFp  *os.File
	index  int32
	errors chan error
	events chan int

	interfaceName string

	// the device was passed IFF_NO_PI
	nopi        bool
	netlinkSock int

	hackListenerClosed sync.Mutex

	// check tun channel was close - signal
	statusListenerShutdown chan struct{}

	// guards calling initNameCache which sets following fields
	closeOnce sync.Once
	nameOnce  sync.Once
	nameError error
}

func NewTUNConfigure(tunName string) *TUNConfigure {
	tfp, err := os.OpenFile(tunDevicePoint, unix.O_RDWR|unix.O_CLOEXEC, 0)
	if err != nil {
		logrus.Fatalf("interface %s create channel error", tunName)
	}

	ifr, err := unix.NewIfreq(tunName)
	if err != nil {
		logrus.Println(err)
		return nil
	}

	ifr.SetUint16(unix.IFF_TUN)
	_, _, errno := unix.Syscall(
		unix.SYS_IOCTL,
		tfp.Fd(),
		uintptr(unix.TUNSETIFF),
		uintptr(unsafe.Pointer(ifr)),
	)

	if errno != 0 {
		logrus.Println(err)
		return nil
	}

	// Set Tun Interface Persistent
	_, _, errno = unix.Syscall(unix.SYS_IOCTL, tfp.Fd(), uintptr(unix.TUNSETPERSIST), uintptr(1))
	if errno != 0 {
		logrus.Println(errno)
		return nil
	}

	return &TUNConfigure{
		tunFp:         tfp,
		interfaceName: tunName,
	}
}

func (tun *TUNConfigure) AssignTunChannelAddress(tunConfig *TUNConfigure, cidr string) error {
	iface, err := netlink.LinkByName(tunConfig.interfaceName)
	if err != nil {
		return err
	}

	addrs, err := netlink.AddrList(iface, netlink.FAMILY_ALL)
	if err != nil {
		return err
	}

	for _, addr := range addrs {
		err = netlink.AddrDel(iface, &addr)
		if err != nil {
			logrus.Println(err)
		}
	}

	addr := &netlink.Addr{
		IPNet: &net.IPNet{
			IP:   net.ParseIP(GetCidrIP(cidr)),
			Mask: net.CIDRMask(31, 32),
		},
	}

	return netlink.AddrAdd(iface, addr)
}

func (tun *TUNConfigure) DeleteTUNChannel(tunConfig *TUNConfigure) {
	iface, err := netlink.LinkByName(tunConfig.interfaceName)
	if err != nil {
		logrus.Warn(err)
	}

	err = netlink.LinkDel(iface)
	if err != nil {
		logrus.Warn(err)
	}
}

func (tun *TUNConfigure) File() *os.File {
	return tun.tunFp
}

func ReIPAddress(fullIP string) string {
	return strings.SplitN(fullIP, ":", 2)[0]
}

func GetCidrIP(cidr string) string {
	_, ipNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return ""
	}

	// 将网络地址转换为无符号整数
	start := ipNet.IP.To4().Mask(ipNet.Mask)

	// 将无符号整数加上偏移量，得到目标 IP 地址
	for i := 1; i < 20; i++ {
		incrementIP(start)
	}

	return start.String()
}

func incrementIP(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}
