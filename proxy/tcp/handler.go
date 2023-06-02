package tcp

import (
	"bufio"
	"context"
	"fmt"
	"github.com/sxueck/ewaf/pkg"
	"github.com/sxueck/ewaf/proxy"
	"io"
	"log"
	"net"
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
	var wg = &sync.WaitGroup{}
	for _, v := range in.([]pkg.Frontend) {
		cv := v
		wg.Add(1)
		go func() {
			defer wg.Done()
			log.Printf("create a new connection tcp channel :%d => %s", cv.ListenPort, (cv.Location)[0].Backend.ByPass)
			lis, err := net.Listen("tcp", fmt.Sprintf(":%d", cv.ListenPort))
			if err != nil {
				log.Fatalf("Failed to start listener: %v", err)
			}

			for {
				client, err := lis.Accept()
				if err != nil {
					log.Printf("Failed to accept connection: %v", err)
					continue
				}

				go handleClient(client, (cv.Location)[0].Backend.ByPass)
			}
		}()
	}

	wg.Wait()
	return nil
}

func handleClient(client net.Conn, targetAddr string) {
	target, err := net.Dial("tcp", targetAddr)
	defer func() {
		client.Close()
		target.Close()
	}()

	tcpConn, ok := target.(*net.TCPConn)
	if ok {
		_ = tcpConn.SetKeepAlive(true)
		_ = tcpConn.SetKeepAlivePeriod(3 * time.Minute)
	}

	if err != nil {
		log.Printf("Failed to connect to target: %v", err)
		return
	}

	reader := bufio.NewReader(client)
	writer := bufio.NewWriter(client)
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
