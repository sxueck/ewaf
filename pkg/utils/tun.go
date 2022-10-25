//go:build linux

package utils

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"os"
	"sync"
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

	return &TUNConfigure{
		tunFp:         tfp,
		interfaceName: tunName,
	}
}

func (tun *TUNConfigure) File() *os.File {
	return tun.tunFp
}
