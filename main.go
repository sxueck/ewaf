package main

import (
	"github.com/sxueck/ewaf/pkg/server"
	"log"
)

func main() {
	var safePort = 8080
	var injectionServer = server.FrontendServer{}
	err := injectionServer.Start(server.Server{
		Name: "test",
		IP:   "127.0.0.1",
		Port: uint8(safePort),
	})
	if err != nil {
		log.Fatal(err)
	}
}
