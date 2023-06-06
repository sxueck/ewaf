package tcp

import (
	"bufio"
	"context"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
	"github.com/sxueck/ewaf/pkg"
	"github.com/sxueck/ewaf/proxy"

	"io"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

type ServerOptions struct {
	ctx    context.Context
	cfg    *pkg.GlobalConfig
	FrMark string // frontend mark
}

func (gso *ServerOptions) WithContext(ctx context.Context, cfg *pkg.GlobalConfig) {
	gso.ctx = ctx
	gso.cfg = cfg
}

func (gso *ServerOptions) Start() any {
	tcpFrs := proxy.CheckTheSurvivalOfUpstreamServices(gso.cfg.Servers, gso.FrMark)
	return tcpFrs
}

func (gso *ServerOptions) Stop() {

}

func (gso *ServerOptions) Serve(in any) error {
	statMap := NewTCPStatMap()
	var wg = &sync.WaitGroup{}

	for _, v := range in.([]pkg.Frontend) {
		cv := v
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("create a new connection tcp channel :%d => %s",
				cv.ListenPort, (cv.Location)[0].Backend.ByPass,
			)

			lis := &CustomRule{}
			lis.IPAddr = net.JoinHostPort(net.IPv4zero.String(), strconv.Itoa(cv.ListenPort))

			err := lis.Listen()
			if err != nil {
				log.Println(err)
				return
			}

			go gso.CaptureTCPPacketFiltering(cv.ListenPort,
				WithTCPServerSYNACKRecv,
			)

			for {
				// 如果不做限制，会出现重复使用已关闭连接的错误
				ctx, cancel := context.WithCancel(context.Background())
				client, err := lis.Accept()
				if err != nil {
					log.Printf("Failed to accept connection: %v", err)
					continue
				}

				go ContinuousGetTCPState(ctx, &client, statMap)
				go handleClient(cancel, &client, (cv.Location)[0].Backend.ByPass)
			}
		}()
	}

	wg.Wait()
	return nil
}

func handleClient(cancel context.CancelFunc, client *net.Conn, targetAddr string) {
	target, err := net.Dial("tcp", targetAddr)
	defer func() {
		(*client).Close()
		target.Close()
		cancel()
	}()

	if (*client).LocalAddr().String() == (*client).RemoteAddr().String() {
		logrus.Warn("Land Attack detected!")
		return
	}

	tcpConn, ok := target.(*net.TCPConn)
	if ok {
		_ = tcpConn.SetKeepAlive(true)
		_ = tcpConn.SetKeepAlivePeriod(3 * time.Minute)
	}

	if err != nil {
		log.Printf("Failed to connect to target: %v", err)
		return
	}

	reader := bufio.NewReader(*client)
	writer := bufio.NewWriter(*client)
	targetReader := bufio.NewReader(target)
	targetWriter := bufio.NewWriter(target)

	go func() {
		_, err = io.Copy(targetWriter, reader)
		if err != nil {
			log.Printf("Error while copying client to target: %v", err)
		}
		_ = targetWriter.Flush()
	}()

	_, err = io.Copy(writer, targetReader)
	if err != nil {
		log.Printf("Error while copying target to client: %v", err)
	}
	_ = writer.Flush()
}

func (gso *ServerOptions) CaptureTCPPacketFiltering(
	port int, opts ...func(*gopacket.PacketSource, chan<- error)) {
	h, err := pcap.OpenLive(
		gso.cfg.Global.Interface,
		1600,
		true,
		pcap.BlockForever,
	)
	if err != nil {
		logrus.Fatal(
			"pre-filter establishment exception, please check whether the network card problems", err)
	}
	defer h.Close()

	filter := fmt.Sprintf("tcp and dst port %d", port)
	err = h.SetBPFFilter(filter)
	if err != nil {
		logrus.Fatal("cannot load filters properly", err)
	}

	pktSource := gopacket.NewPacketSource(h, h.LinkType())
	var e = make(chan error, 1)
	for _, v := range opts {
		go v(pktSource, e)
	}
	err = <-e
	if err != nil {
		logrus.Fatal("filtering exception", err)
	}
}
