package main

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"net"
)

func main() {
	srcIP := net.IPv4(127, 0, 0, 1)
	dstIP := net.IPv4(172, 26, 247, 15)
	device := "eth0"
	CreatePacket(srcIP, dstIP, device)
}

func CreatePacket(src, dst net.IP, device string) {
	handle, err := pcap.OpenLive(device, 65536, false, pcap.BlockForever)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()

	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	// 创建并设置以太网层
	ethernetLayer := &layers.Ethernet{
		SrcMAC:       net.HardwareAddr{0xFF, 0xAA, 0xFA, 0xAA, 0xFF, 0xAA},
		DstMAC:       net.HardwareAddr{0xBD, 0xBD, 0xBD, 0xBD, 0xBD, 0xBD},
		EthernetType: layers.EthernetTypeIPv4,
	}

	// 创建并设置 IPv4 层
	ipLayer := &layers.IPv4{
		Version:  4,
		TTL:      64,
		SrcIP:    src,
		DstIP:    dst,
		Protocol: layers.IPProtocolTCP,
	}

	// 创建并设置 TCP 层
	tcpLayer := &layers.TCP{
		SrcPort: 4321,
		DstPort: 8089,
		Seq:     11050,
		SYN:     true,
		ACK:     true,
		Window:  14600,
	}

	// 将网络层关联到 TCP 层以计算校验和
	err = tcpLayer.SetNetworkLayerForChecksum(ipLayer)
	if err != nil {
		log.Fatal(err)
	}

	// 将以太网层、IPv4 层和 TCP 层序列化到缓冲区
	buffer := gopacket.NewSerializeBuffer()
	err = gopacket.SerializeLayers(buffer, options, ethernetLayer, ipLayer, tcpLayer)
	if err != nil {
		log.Fatal(err)
	}

	// 获取序列化后的数据并发送
	outgoingPacket := buffer.Bytes()
	err = handle.WritePacketData(outgoingPacket)
	if err != nil {
		log.Fatal(err)
	}
}
