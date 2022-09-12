package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"text/tabwriter"

	"github.com/miekg/dns"
)

// 遍历字典，枚举子域，并且查询子域名的a记录和cname记录
// 使用 go run subdomain_gusee.go -domain="xxx.xx" -wordlist="xx" -c=10 -dns="8.8.8.8:53"

var (
	fdomain    = flag.String("domain", "baidu.com", "猜解的域")
	fwordlist  = flag.String("wordlist", "", "暴破字典")
	fworkcount = flag.Int("c", 30, "协程数")
	fdns       = flag.String("dns", "8.8.8.8:53", "dns")
)

type result struct {
	ip      string // 解析的a记录
	host    string
	allHost []string // 所有的cname  可能cname1->cname2->cname3->a
}

// 查询a记录
func lookupA(fqdn, fdns string) ([]string, error) {
	var m dns.Msg
	var ips []string
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeA)
	in, err := dns.Exchange(&m, fdns)

	if err != nil {
		return ips, err
	}

	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}
	return ips, nil
}

// 查询a记录
func lookupCname(fqdn, fdns string) ([]string, error) {
	var m dns.Msg
	var fqdns []string
	m.SetQuestion(dns.Fqdn(fqdn), dns.TypeCNAME)
	in, err := dns.Exchange(&m, fdns)

	if err != nil {
		return fqdns, err
	}

	for _, answer := range in.Answer {
		if a, ok := answer.(*dns.CNAME); ok {
			log.Println(a.Target)
			fqdns = append(fqdns, a.Target)
		}
	}
	return fqdns, nil
}

func lookup(fqdn, fdns string) []result {
	var results []result
	var cnames []string
	var cfqdn = fqdn //拷贝出一份 防止篡改
	for {
		cname, err := lookupCname(cfqdn, fdns)
		if err != nil {
			//log.Println(err)
			return nil
		}
		if len(cname) > 0 {
			cfqdn = cname[0] // 跟踪出所有的cname，直到最后一个
			cnames = append(cnames, cname...)
			continue
		}
		ips, err := lookupA(cfqdn, fdns)
		if err != nil {
			//log.Println(err)
			return nil
		}

		for _, v := range ips {
			results = append(results, result{
				ip:      v,
				host:    cfqdn,
				allHost: cnames,
			})
		}
	}
	fmt.Println(results)
	return results
}

// 协程处理
// in输入fqdn  out输出result，在main中接收，isdone 标识是否结束
func worker(in chan string, out chan []result, isdone chan struct{}, fdns string) {
	for item := range in { // 从in队列中读取
		results := lookup(item, fdns)
		if len(results) > 0 {
			out <- results // 将result传出去
		}
	}
	isdone <- struct{}{} // 防止没有完成就终止
}

func main() {
	flag.Parse()

	if *fdomain == "" || *fwordlist == "" {
		fmt.Println("传错错误")
		os.Exit(1)
	}
	in := make(chan string)
	out := make(chan []result)
	isdone := make(chan struct{})

	// 暴破文件读取
	fn, err := os.Open(*fwordlist) // 打开文件
	if err != nil {
		log.Fatalln(err)
	}
	defer fn.Close()
	// 创建遍历器
	scanner := bufio.NewScanner(fn)

	// 创建10个消费者
	for i := 0; i < *fworkcount; i++ {
		go worker(in, out, isdone, *fdns)
	}

	// 需要先创建等待任务的协程，不然会死锁
	// 生产者
	for scanner.Scan() {
		in <- fmt.Sprintf("%s.%s", scanner.Text(), *fdomain)
	}

	// 结果
	var results []result
	go func() {
		for item := range out {
			results = append(results, item...)
		}
		isdone <- struct{}{}
	}()

	close(in) // 所有数据发送完成，进行关闭管道

	// 回收worker进程的完成管道
	for i := 0; i < *fworkcount; i++ {
		<-isdone
	}
	// 接收完worker的数据 关闭
	close(out)

	<-isdone // result的is_done 管道

	w := tabwriter.NewWriter(os.Stdout,0,8,' ',' ',0) //NewWriter(output io.Writer, minwidth int, tabwidth int, padding int, padchar byte, flags uint) *tabwriter.Writer
	for _,v := range results{
		fmt.Fprintf(w,"%s\t%s\n",v.host,v.ip)
	}
	w.Flush()

}
