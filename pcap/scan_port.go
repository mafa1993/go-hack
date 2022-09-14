package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

// 端口扫描 查看端口是否有其他数据包，如果有表明端口真的开放

var (
	snaplen  = int32(1600) // 每帧捕获的数据量
	promisc  = true        // 是否开启混杂模式
	timeout  = pcap.BlockForever
	filter   = "tcp[13]==0x11 or tcp[13] == 0x10 or tcp[13]==0x18" // 设置过滤器，检查tcp头的第14个字节
	devFound = false
	results  = make(map[string]int)
)

func main() {
	if len(os.Args) != 4 {
		log.Fatalln("参数错误")
	}

	devices, err := pcap.FindAllDevs() // 获取所有设备
	if err != nil {
		log.Panic(err)
	}

	iface := os.Args[1]
	// 网卡名验证
	for _, dev := range devices {

		if dev.Name == iface {
			devFound = true
		}
	}

	if !devFound {
		log.Fatalln("没有找到指定的网卡")
	}

	ip := os.Args[2]
	go capture(iface, ip)
	time.Sleep(time.Second)

	ports, _ := explode(os.Args[3])

	for _, port := range ports {
		target := fmt.Sprintf("%s,%s", ip, port)
		fmt.Println("trying ", target)
		c, err := net.DialTimeout("tcp", target, 1000*time.Millisecond)

		if err != nil {
			continue
		}

		c.Close()
	}
	time.Sleep(time.Second)

	for port, confidence := range results {
		fmt.Printf("端口%s %d", port, confidence)
	}

}

func capture(iface, target string) {
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
		networkLayer := packet.NetworkLayer() // 获取网络层数据
		if networkLayer == nil {
			continue
		}

		transportLayer := packet.TransportLayer() // 获取传输层
		if transportLayer == nil {
			continue
		}

		srcHost := networkLayer.NetworkFlow().Src().String() // 获取源ip
		if srcHost != target {
			log.Fatalln("srchost " + srcHost)
			continue
		}
		srcPort := transportLayer.TransportFlow().Src().String()
		results[srcPort] += 1
	}
}

func explode(portString string) ([]string, error) {
	ret := make([]string, 0)

	ports := strings.Split(portString, ",")
	for _, port := range ports {
		port := strings.TrimSpace(port)
		ret = append(ret, port)
	}

	return ret, nil
}
