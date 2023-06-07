package tcp

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
	"io"
	"net"
	"os"
	"sync"
	"time"
)

// 使用splice的时候，需要将文件描述符设置为非阻塞模式
func setNonblock(conn net.Conn) (*os.File, error) {
	tcpConn := conn.(*net.TCPConn)
	fp, err := tcpConn.File()
	if err != nil {
		return nil, err
	}

	if err = unix.SetNonblock(int(fp.Fd()), true); err != nil {
		return nil, err
	}

	return fp, nil
}

func splicePipes(src, dst *os.File) error {
	bufSize := os.Getpagesize()
	var written int64

	r, w, err := os.Pipe()
	if err != nil {
		return err
	}

	defer func() {
		r.Close()
		w.Close()
	}()

	for {
		nr, err := unix.Splice(int(src.Fd()), nil, int(w.Fd()), nil, bufSize, 0)
		if err != nil {
			if err == unix.EAGAIN {
				continue
			}

			if err == unix.EINVAL {
				// 不支持splice则回退到io.copy
				_, err := io.Copy(dst, src)
				if err != nil {
					return fmt.Errorf("io.copy error: %v", err)
				}
			}
			return err
		}

		if nr > 0 {
			_, err := unix.Splice(int(r.Fd()), nil, int(dst.Fd()), nil, int(nr), 0)
			if err != nil {
				return err
			}
		}

		written += nr
		if nr == 0 {
			break
		}
	}

	return nil
}

func handleClient(cancel context.CancelFunc, cfd *os.File, targetAddr string) {
	target, err := net.Dial("tcp", targetAddr)
	defer cancel()

	tcpConn, ok := target.(*net.TCPConn)
	if ok {
		_ = tcpConn.SetKeepAlive(true)
		_ = tcpConn.SetKeepAlivePeriod(3 * time.Minute)
		_ = tcpConn.SetNoDelay(true)
	}

	if err != nil {
		logrus.Printf("Failed to connect to target: %v", err)
		return
	}

	tfd, err := setNonblock(target)
	if err != nil {
		logrus.Printf("Failed to set target to nonblock: %v", err)
		return
	}

	wg := sync.WaitGroup{}

	go func() {
		wg.Add(1)
		defer wg.Done()
		err = splicePipes(cfd, tfd)
		if err != nil {
			logrus.Printf("error while copying client to target: %v", err)
		}
	}()

	go func() {
		wg.Add(1)
		defer wg.Done()
		err = splicePipes(tfd, cfd)
		if err != nil {
			logrus.Printf("error while copying target to client: %v", err)
		}
	}()

	wg.Wait()
}
