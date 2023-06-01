package tcp

import (
	"context"
	"fmt"
	"github.com/sxueck/ewaf/pkg"
	"github.com/sxueck/ewaf/proxy"
	"io"
	"log"
	"net"
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

func (gso *ServerOptions) Start() error {
	tcpFrs := proxy.CheckTheSurvivalOfUpstreamServices(*gso.cfg.Servers, gso.FrMark)

	log.Printf("%+v\n", tcpFrs)
	for _, v := range tcpFrs {
		cv := v
		go func() {
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

				go handleClient(client, (*cv.Location)[0].Backend.ByPass)
			}
		}()
	}

	return nil
}

func (gso *ServerOptions) Stop() {

}

func (gso *ServerOptions) Serve() error {
	return nil
}

func handleClient(client net.Conn, targetAddr string) {
	target, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Failed to connect to target: %v", err)
		client.Close()
		return
	}

	go func() {
		_, err := io.Copy(target, client)
		if err != nil {
			log.Printf("Error while copying client to target: %v", err)
		}
	}()

	_, err = io.Copy(client, target)
	if err != nil {
		log.Printf("Error while copying target to client: %v", err)
	}

	client.Close()
	target.Close()
}
