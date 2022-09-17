package main

import (
	"C"
	"fmt"
)

// 下面的注释代表导出一个函数叫Start

//export Start
func Start() {
	fmt.Println("go")
}

// main为空即可
func main() {}
