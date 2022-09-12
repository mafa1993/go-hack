package main

import (
	"flag"
	"fmt"

	"github.com/miekg/dns"
)

// dns a记录查询
//用法 go run .\search_a.go -host="www.taobao.com" -dns="8.8.8.8:53" 

var (
	host     = flag.String("host", "www.baidu.com", "what host you want to search")
	dns_host = flag.String("dns", "114.114.114.114:53", "dns host you want to use")
)

func main() {
	flag.Parse()
	var msg dns.Msg                  // dns查询返回的数据
	fqdn := dns.Fqdn(*host)          // 将host转换为Fully Qualified Domain Name
	msg.SetQuestion(fqdn, dns.TypeA) // 设置查询a记录
	in, err := dns.Exchange(&msg, *dns_host)
	if err != nil {
		panic(err)
	}
	for _, item := range in.Answer {
		if a, ok := item.(*dns.A); ok { // 进行断言，如果是*dns.A才打印
			fmt.Println(a.A) // a.A为net.IP类型，但是实现了Stirng方法
		}
	}
}
