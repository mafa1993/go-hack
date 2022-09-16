使用syscall可以加载驻留内存的恶意软件或者钩子函数

c++和go类型

c++ | go 
---|---
boolean|byte
bool|int32
byte|byte
dword|uint32
dword32|uint32
dword64|uint46
word|uint16
handle|uintptr
lpvoid | uintptr
size_t | uintptr
lpcvoid | uintptr
hmodule | uintptr
lpcstr | uintptr
lpdword | uintptr

## 进程注入的流程
 
1. 使用openProcess()建立进程句柄和进程访问权限
2. VirtualAllocEx() 分配虚拟内存，
3. WriteProcessMemory()将shellcode或者dll加载到进程内存中
4. 使用CreateRemoteThread()调用本地导出的dll函数，使得第三步写入内存的字节码执行


# 扩展

1. 可以使用Process Hacker Process Monitor工具查看进程状态
2. 可以不通过dllpath加载dll   使用cs或者msfvenom生成shellcode
3. 可以将frida加载进去，执行js

# 使用

go run .\main.go .\helper.go .\inject.go .\tokens.go -pid=10676 -dll="C:\Windows\System32.dll"

pid为10676dll 为C:\Windows\System32.dllprocess handle0x158
申请内存 0x1ae387b0000
kernal dll 的地址为0x7ffd66100000load memory 0x7ffd6611ebb0
Thread 创建0xc0000a60c8
thread create 0x160
