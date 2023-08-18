package tcp

import (
	"context"
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
	"github.com/sirupsen/logrus"
	"github.com/sxueck/ewaf/pkg"
	"github.com/sxueck/ewaf/pkg/infra"
	"github.com/sxueck/ewaf/proxy"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ServerOptions struct {
	ctx    context.Context
	cfg    *pkg.GlobalConfig
	FrMark string // frontend mark
	bloom  *infra.Filter
}

func (gso *ServerOptions) WithContext(ctx context.Context, cfg *pkg.GlobalConfig) {
	gso.ctx = ctx
	gso.cfg = cfg
}

func (gso *ServerOptions) ExtraFrMark() string {
	return gso.FrMark
}

func (gso *ServerOptions) Start() any {
	tcpFrs := proxy.CheckTheSurvivalOfUpstreamServices(gso.cfg.Servers, gso.FrMark)
	gso.bloom = infra.NewBloom(2048, 3, false)

	go func() {
		// 定期重置布隆过滤器，动态拒绝连接
		ticker := time.NewTicker(1 * time.Microsecond)
		defer ticker.Stop()

		for {
			<-ticker.C
			gso.bloom.Reset()
			LoadDenyIPRules(gso, nil)
			ticker.Reset(NextStatusInterval * time.Second * 60) // 120s
		}
	}()
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

			// 过滤规则加载
			go gso.CaptureTCPPacketFiltering(cv.ListenPort,
				WithTCPServerSYNACKRecv,
			)

			for {
				// 如果不做限制，会出现重复使用已关闭连接的错误
				ctx, cancel := context.WithCancel(context.Background())
				f, err := lis.AcceptToFD()
				if err != nil {
					log.Printf("Failed to accept connection: %v", err)
					continue
				}

				client, _ := net.FileConn(f)
				if gso.bloom.KeySize() > 0 {
					dip := client.RemoteAddr().String()
					if gso.bloom.TestString(ReIPAddress(dip)) {
						logrus.Println("deny ip:", dip)
						_ = client.Close()
						continue
					}
				}

				go ContinuousGetTCPState(ctx, &client, statMap)
				// f为请求端的socket描述符
				go handleClient(cancel, f, (cv.Location)[0].Backend.ByPass)
			}
		}()
	}

	wg.Wait()
	return nil
}

func (gso *ServerOptions) CaptureTCPPacketFiltering(
	port int, opts ...func(<-chan gopacket.Packet, *ServerOptions)) {

	iface, err := net.InterfaceByName(gso.cfg.Global.Interface)
	if err != nil {
		log.Fatal(err)
	}
	h, err := pcap.OpenLive(
		iface.Name, 1600, true,
		pcap.BlockForever,
	)
	if err != nil {
		logrus.Fatal(
			"pre-filter establishment exception, please check whether the network card problems", err)
	}

	filter := fmt.Sprintf("tcp and dst port %d", port)
	err = h.SetBPFFilter(filter)
	if err != nil {
		logrus.Fatal("cannot load filters properly", err)
	}

	pktSource := gopacket.NewPacketSource(h, h.LinkType())
	for _, v := range opts {
		go v(pktSource.Packets(), gso)
	}
}

func ReIPAddress(fullIP string) string {
	return strings.SplitN(fullIP, ":", 2)[0]
}
