package main

import (
	"fmt"
	"net"
)

func main() {
	targetIP := "172.26.247.15" // 目标IP地址
	targetPort := 8089          // 目标端口号

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", targetIP, targetPort))
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	// 发送 SYN+ACK 包
	_, err = conn.Write([]byte{0x12, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Println("SYN+ACK packet sent successfully!")
}
