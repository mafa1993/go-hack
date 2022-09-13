package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/miekg/dns"
)

// 检查dns消息，并进行适当的路由
// 功能1. 创建处理函数接收传入的查询
// 2. 检查dns question 并提取域名
// 3. 识别与域名相关的dns服务器
// 4. 向上游dns服务器发送请求，并返回

// 读取配置文件
func parse(filename string) (map[string]string, error) {
	records := make(map[string]string)

	fn, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(fn)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, ",", 2)
		if len(parts) < 2 {
			return nil, errors.New("文件格式不对")
		}

		records[parts[0]] = parts[1]
	}
	return records, err
}

func main() {
	var lock sync.RWMutex
	records, err := parse("proxy.conf")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", records)
	// 设置dns服务
	dns.HandleFunc(".", func(w dns.ResponseWriter, m *dns.Msg) {
		// 没有请求直接返回
		if len(m.Question) == 0 {
			dns.HandleFailed(w, m)
			return
		}

		fqdn := m.Question[0].Name //获取请求查询的域名
		var parts []string
		parts = strings.Split(fqdn, ".")
		if len(parts) > 2 {
			fqdn = strings.Join(parts[len(parts)-2:], ".")
		}
		lock.RLock()           // 加写锁，保证读的时候不能写
		match := records[fqdn] // 获取域名对应的dns服务器
		lock.Unlock()

		// 没有匹配上
		if match == "" {
			dns.HandleFailed(w, m)
			return
		}

		resp, err := dns.Exchange(m, match) // 给上游发送请求
		if err != nil {
			dns.HandleFailed(w, m)
			return
		}

		// 将上游返回写给客户端
		if err := w.WriteMsg(resp); err != nil {
			dns.HandleFailed(w, m)
			return
		}

	})

	// 信号量处理，重新加载配置
	go func() {
		sigs := make(chan os.Signal, 1) // 如果构建成无缓冲的，还有开个协程
		signal.Notify(sigs, syscall.SIGUSR1)
		for s := range sigs {
			switch s {
			case syscall.SIGUSR1:
				fmt.Println("sig接收")
				lock.Lock() // 重新加载配置文件
				parse("proxy.conf")
				lock.Unlock()
			}
		}

	}()

	log.Fatal(dns.ListenAndServe(":53", "udp", nil))
}
