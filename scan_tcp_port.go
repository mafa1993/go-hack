package main

import (
	"fmt"
	"net"
)

func main() {
	var datas []int // 用来保存扫描到的port
	var ports chan int = make(chan int)
	var done chan int = make(chan int)

	// 创建消费者
	for i := 0; i < 10; i++ {
		go func() {
			for {
				port, ok := <-ports
				if !ok {
					break
				}
				add := fmt.Sprintf("scanme.namp.com:%d", port)
				_, err := net.Dial("tcp", add)
				if err != nil {
					//log.Fatalln(err)
					fmt.Println(err)
				} else {
					datas = append(datas, port)
				}

			}
			done <- 1
		}()
	}

	for i := 1; i < 100; i++ {
		ports <- i
	}

	for v := range done {
		fmt.Println(v)
	}
	fmt.Println(datas)

}
