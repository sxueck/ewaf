//go:build linux

package water

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
	"strings"
	"sync"
	"unsafe"
)

const (
	defaultMTU     = 1420
	ifReqSize      = unix.IFNAMSIZ + 0x40
	tunDevicePoint = "/dev/net/tun"
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

func (tun *TUNConfigure) File() *os.File {
	return tun.tunFp
}

func ReIPAddress(fullIP string) string {
	return strings.SplitN(fullIP, ":", 2)[0]
}
