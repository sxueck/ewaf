package main

import (
	"github.com/sxueck/ewaf/pkg/server"
	"log"
)

func main() {
	var injectionServer = server.FrontendServer{}
	err := injectionServer.Start(server.Server{
		Name: "test",
		IP: "127.0.0.1",
		Port: 8080,
	})
	if err != nil {
		log.Fatal(err)
	}
}
