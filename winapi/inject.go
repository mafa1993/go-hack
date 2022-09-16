package main

import (
	"fmt"
	"syscall"
	"unsafe"

	"github.com/pkg/errors"
)

var nullRef int

// 获取进程句柄
func OpenProcessHandle(i *Inject) error {
	// 定义进程访问权限
	var right uint32 = PROCESS_CREATE_THREAD |
		PROCESS_QUERY_INFORMATION | // 查询进程信息
		PROCESS_VM_OPERATION |
		PROCESS_VM_READ |
		PROCESS_VM_WRITE

	var inheritHanle uint = 0 // 新进程句柄是否继承现有句柄
	var processId uint32 = i.Pid

	// 获取进程句柄
	remoteProcHandle, _, lastErr := ProcOpenProcess.Call( // 进行系统调用
		uintptr(right),        //DWORD dwDesiredAccess
		uintptr(inheritHanle), // Bool bInheritHandle
		uintptr(processId),    //DWORD dwProcessId
	)
	if remoteProcHandle == 0 {
		return errors.Wrap(lastErr, "不能获取句柄")
	}

	i.RemoteProcHandle = remoteProcHandle // 记录到结构体
	fmt.Printf("pid为%v", i.Pid)
	fmt.Printf("dll 为%v", i.DllPath)
	fmt.Printf("process handle%v \n", unsafe.Pointer(i.RemoteProcHandle))
	return nil
}

// 内存申请
func VirtualAllocEx(i *Inject) error {
	var flAllocationType uint32 = MEM_COMMIT | MEM_RESERVE
	var flProtect uint32 = PAGE_EXECUTE_READWRITE

	lpBaseAddress, _, lastErr := ProcVirtualAllocEx.Call(
		i.RemoteProcHandle, // HANDLE hProcess
		uintptr(nullRef),   // LPVOID ipadress
		uintptr(i.DLLSize), // size_t
		uintptr(flAllocationType),
		uintptr(flProtect), //dword flprotect
	)

	if lpBaseAddress == 0 {
		return errors.Wrap(lastErr, "申请内存失败")
	}

	i.Lpaddr = lpBaseAddress
	fmt.Printf("申请内存 %v\n", unsafe.Pointer(i.Lpaddr))
	return nil
}

func WriteProcessMemory(i *Inject) error {
	var nBytesWritten *byte
	dllPathBytes, err := syscall.BytePtrFromString(i.DllPath) // 接收一个dll地址，返回生成的字节切片的地址

	if err != nil {
		return err
	}

	writeMem, _, lastErr := ProcWriteProcessMemory.Call(
		i.RemoteProcHandle,                     // HANDLE hprocess
		i.Lpaddr,                               // LPVOID lpbaseAddress
		uintptr(unsafe.Pointer(dllPathBytes)),  // LPCVOID
		uintptr(i.DLLSize),                     // size_t
		uintptr(unsafe.Pointer(nBytesWritten)), //size_t *lpNumberOfBytesWriten
	)

	if writeMem == 0 {
		return errors.Wrap(lastErr, "shellcode 写入内存失败")
	}

	return nil
}

// loadlibrarya 会将指定的模块加载到进程调用的内存空间中，所以需要得到library的内存位置
func GetLoadLibAddress(i *Inject) error {
	var llibBytesPtr *byte
	llibBytesPtr, err := syscall.BytePtrFromString("LoadLibraryA") // 返回这个字符串在内存中的位置

	if err != nil {
		return err
	}

	lladdr, _, lastErr := ProcGetProcAddress.Call(
		ModKernel32.Handle(),
		uintptr(unsafe.Pointer(llibBytesPtr)), // LPCSTR lpProcname
	)

	if &lladdr == nil {
		return errors.Wrap(lastErr, "没有找到地址")
	}
	i.LoadLibAddr = lladdr
	fmt.Printf("kernal dll 的地址为%v", unsafe.Pointer(ModKernel32.Handle()))

	fmt.Printf("load memory %v\n", unsafe.Pointer(i.LoadLibAddr))

	return nil
}

// 针对远程进程的虚拟内存区域创建一个线程
func CreateRemoteThread(i *Inject) error {
	var threadId uint32 = 0
	var dwCreateionFlags uint32 = 0

	remoteThread, _, lastErr := ProcCreateRemoteThread.Call(
		i.RemoteProcHandle,
		uintptr(nullRef),
		uintptr(nullRef),
		i.LoadLibAddr,
		i.Lpaddr, // 虚拟内存位置
		uintptr(dwCreateionFlags),
		uintptr(unsafe.Pointer(&threadId)),
	)

	if remoteThread == 0 {
		return errors.Wrap(lastErr, "创建线程失败")
	}

	i.RThread = remoteThread
	fmt.Printf("Thread 创建%v\n", unsafe.Pointer(&threadId))
	fmt.Printf("thread create %v\n", unsafe.Pointer(i.RThread))

	return nil
}

// 识别特定对象合适处于发信号的状态
func WaitForSingleObject(i *Inject) error {
	var dwMilliseconds uint32 = INFINITE
	var dwExitCode uint32
	rWaitValue, _, lastErr := ProcWaitForSingleObject.Call(
		i.RThread,
		uintptr(dwMilliseconds),
	)

	if rWaitValue != 0 {
		return errors.Wrap(lastErr, "线程状态错误")
	}

	success, _, lastErr := ProcGetExitCodeThread.Call(
		i.RThread,
		uintptr(unsafe.Pointer(&dwExitCode)),
	)

	if success == 0 {
		return errors.Wrap(lastErr, "退出码不对")
	}

	closed, _, lastErr := ProcCloseHandle.Call(i.RThread)

	if closed == 0 {
		return errors.Wrap(lastErr, "关闭错误")
	}

	return nil
}

func VirtualFreeEx(i *Inject) error {
	var dwFreeType uint32 = MEM_RELEASE
	var size uint32 = 0 //Size must be 0 if MEM_RELEASE all of the region
	rFreeValue, _, lastErr := ProcVirtualFreeEx.Call(
		i.RemoteProcHandle,
		i.Lpaddr,
		uintptr(size),
		uintptr(dwFreeType))
	if rFreeValue == 0 {
		return errors.Wrap(lastErr, "释放内存出错")
	}
	fmt.Println("释放内存成功")
	return nil
}
