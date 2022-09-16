package main

import (
	"fmt"
	"runtime"
	"unsafe"
)

// 指针相关

// uintptr 和 unsafe.Pointer
// 1. unsafe.Pointer和任何指针可以互相转换
// 2. uintptr和unsafe.Pointer可以互相转换
// 3. uintprt本质上是一个无符号整箱，如果使用unsafe.Pointer创建一个指针赋给unintptr,go中无法保存引用关系

func main() {
	state()
}

func state() {
	var onload = createPtr("onload")
	var success = createPtr("success")
	var receive = createPtr("receive")

	maps := make(map[string]interface{})
	maps["onload"] = unsafe.Pointer(onload)
	maps["success"] = unsafe.Pointer(success)
	maps["receive"] = uintptr(unsafe.Pointer(receive)) // uintptr 是一个整型，这里执行完，receive会被释放，可能不会找到真正的receive值

	fmt.Println(*(*string)(maps["onload"].(unsafe.Pointer))) // 安全，会打印出原始值

	fmt.Println(*(*string)(unsafe.Pointer(maps["receive"].(uintptr)))) // success被回收  找不到success值

	runtime.KeepAlive(receive) //显示告诉go，手工释放receive，或者在运行结束时，扔保持receive可以访问，防止被回收
}

func createPtr(s string) *string {
	return &s
}
