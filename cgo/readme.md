
# 准备工作

1. 安装msys2 对c进行编译

# go里面嵌入c

1. go支持在编译时对注释进行编译，支持多种语言的注释编译
2. cgo 可以go代码调用c或者c调用go，如果先要使用go的垃圾回收和内存管理，需要在go中申请内存，传给c，除非使用free(), 不然不会对内存进行释放，使用defer可以确保go中引用的所有c内存都会进行垃圾回收

# c调go

1. 将go编译到一个归档文件中，然后c将归档文件编译到dll中
2. 需要先将go编译成插件  go build -buildmode=c-archive   会生成一个.a文件
3.  gcc -lpthread .\main.c .\main.a  生成exe
4.  gcc -shared -pthread -o x.dll .\main.c .\main.a  生成dll


```
go-hack\cgo\shellcode> .\a.exe
go
```

使用sRDI(git上搜索)可以将dll转化为shellcode

使用PowerSploit实现shellcode注入