package main

import (
	"debug/pe"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	f, err := os.Open("G:\\Windows\\WinSxS\\amd64_microsoft-windows-calc_31bf3856ad364e35_10.0.18362.1_none_7c1b713697f466dd\\calc.exe") // Modify for binary or change to accept args
	check(err)
	pefile, err := pe.NewFile(f) // 创建pe文件对象
	check(err)
	defer f.Close()
	defer pefile.Close()

	dosHeader := make([]byte, 96) // 读取前96个字节
	sizeOffset := make([]byte, 4)

	// Dec to Ascii (searching for MZ)
	_, err = f.Read(dosHeader)
	check(err)
	fmt.Println("[-----DOS Header / Stub----- header和stub解析]")
	fmt.Printf("[+] Magic Value: %s%s\n", string(dosHeader[0]), string(dosHeader[1]))

	// Validate PE+0+0 (Valid PE format)
	pe_sig_offset := int64(binary.LittleEndian.Uint32(dosHeader[0x3c:])) // 从0x3c后开始读取
	f.ReadAt(sizeOffset, pe_sig_offset)                                  // sizeoffset为buffer，pe_sig_offset是读取的位置  这里不加:可以吗？？ 读取0x50 0x45 0x00 0x00

	fmt.Println("[-----Signature Header-----]")
	fmt.Printf("[+] LFANEW Value: %s\n", string(sizeOffset))

	// Create the reader and read COFF Header
	sr := io.NewSectionReader(f, 0, 1<<63-1)       // 读取到2^64次方-1
	_, err = sr.Seek(pe_sig_offset+4, os.SEEK_SET) // 重置指针
	check(err)
	binary.Read(sr, binary.LittleEndian, &pefile.FileHeader) // 二进制读取

	// Get size of OptionalHeader
	// 可选头解析  向加载程序提供重要数据，加载程序将可执行文件加载到虚拟内存
	var sizeofOptionalHeader32 = uint16(binary.Size(pe.OptionalHeader32{}))
	var sizeofOptionalHeader64 = uint16(binary.Size(pe.OptionalHeader64{}))
	var oh32 pe.OptionalHeader32
	var oh64 pe.OptionalHeader64

	// type FileHeader struct {
	// 	Machine              uint16
	// 	NumberOfSections     uint16
	// 	TimeDateStamp        uint32
	// 	PointerToSymbolTable uint32
	// 	NumberOfSymbols      uint32
	// 	SizeOfOptionalHeader uint16
	// 	Characteristics      uint16
	// }

	// Read OptionalHeader
	switch pefile.FileHeader.SizeOfOptionalHeader {
	case sizeofOptionalHeader32:
		binary.Read(sr, binary.LittleEndian, &oh32)
	case sizeofOptionalHeader64:
		binary.Read(sr, binary.LittleEndian, &oh64)
	}
	// Print File Header
	fmt.Println("[-----COFF File Header-----]")
	fmt.Printf("[+] Machine Architecture: %#x\n", pefile.FileHeader.Machine)
	fmt.Printf("[+] Number of Sections: %#x\n", pefile.FileHeader.NumberOfSections)
	fmt.Printf("[+] Size of Optional Header: %#x\n", pefile.FileHeader.SizeOfOptionalHeader)
	// Print section names
	fmt.Println("[-----Section Offsets-----]")
	fmt.Printf("[+] Number of Sections Field Offset: %#x\n", pe_sig_offset+6)
	// this is the end of the Signature header (0x7c) + coff (20bytes) + oh32 (224bytes)
	fmt.Printf("[+] Section Table Offset: %#x\n", pe_sig_offset+0xF8)
	// Print Optional Header
	fmt.Println("[-----Optional Header-----]")

	// 如果对pe文件加入后门，需要对下面两项有所了解
	fmt.Printf("[+] Entry Point: %#x\n", oh32.AddressOfEntryPoint) // 相对于imageBase的可执行代码位置
	fmt.Printf("[+] ImageBase: %#x\n", oh32.ImageBase)             // 将图像加载到内存时第一个字节位置

	fmt.Printf("[+] Size of Image: %#x\n", oh32.SizeOfImage) // 图像的实际大小
	fmt.Printf("[+] Sections Alignment: %#x\n", oh32.SectionAlignment)
	fmt.Printf("[+] File Alignment: %#x\n", oh32.FileAlignment)
	fmt.Printf("[+] Characteristics: %#x\n", pefile.FileHeader.Characteristics)
	fmt.Printf("[+] Size of Headers: %#x\n", oh32.SizeOfHeaders)
	fmt.Printf("[+] Checksum: %#x\n", oh32.CheckSum)
	fmt.Printf("[+] Machine: %#x\n", pefile.FileHeader.Machine)
	fmt.Printf("[+] Subsystem: %#x\n", oh32.Subsystem)
	fmt.Printf("[+] DLLCharacteristics: %#x\n", oh32.DllCharacteristics)

	// Print Data Directory
	fmt.Println("[-----Data Directory----- 数据目录解析，可选头的最后128字节]")
	var winnt_datadirs = []string{
		"IMAGE_DIRECTORY_ENTRY_EXPORT",
		"IMAGE_DIRECTORY_ENTRY_IMPORT",
		"IMAGE_DIRECTORY_ENTRY_RESOURCE",
		"IMAGE_DIRECTORY_ENTRY_EXCEPTION",
		"IMAGE_DIRECTORY_ENTRY_SECURITY",
		"IMAGE_DIRECTORY_ENTRY_BASERELOC",
		"IMAGE_DIRECTORY_ENTRY_DEBUG",
		"IMAGE_DIRECTORY_ENTRY_COPYRIGHT",
		"IMAGE_DIRECTORY_ENTRY_GLOBALPTR",
		"IMAGE_DIRECTORY_ENTRY_TLS",
		"IMAGE_DIRECTORY_ENTRY_LOAD_CONFIG",
		"IMAGE_DIRECTORY_ENTRY_BOUND_IMPORT",
		"IMAGE_DIRECTORY_ENTRY_IAT",
		"IMAGE_DIRECTORY_ENTRY_DELAY_IMPORT",
		"IMAGE_DIRECTORY_ENTRY_COM_DESCRIPTOR",
		"IMAGE_NUMBEROF_DIRECTORY_ENTRIES",
	}
	for idx, directory := range oh32.DataDirectory {
		fmt.Printf("[!] Data Directory: %s\n", winnt_datadirs[idx])
		fmt.Printf("[+] Image Virtual Address: %#x\n", directory.VirtualAddress)
		fmt.Printf("[+] Image Size: %#x\n", directory.Size)
	}

	fmt.Println("[-----Section Table----- 解析分区表]")
	// 包含了相关分区的详细信息，与coff file header中numberofsections 匹配
	for _, section := range pefile.Sections {
		fmt.Println("[+] --------------------")
		fmt.Printf("[+] Section Name: %s\n", section.Name)
		fmt.Printf("[+] Section Characteristics: %#x\n", section.Characteristics)
		fmt.Printf("[+] Section Virtual Size: %#x\n", section.VirtualSize)
		fmt.Printf("[+] Section Virtual Offset: %#x\n", section.VirtualAddress)
		fmt.Printf("[+] Section Raw Size: %#x\n", section.Size)
		fmt.Printf("[+] Section Raw Offset to Data: %#x\n", section.Offset)
		fmt.Printf("[+] Section Append Offset (Next Section): %#x\n", section.Offset+section.Size)
	}

	// s := pefile.Section(".text")
	// fmt.Printf("%v", *s)

	// "Section Table Offset" + (40bytes * number of sections)

}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// Things to consider when modifying an image
// 1. Round off the last section
// 2. Add the stub in the PE Header
// 3. Increase the number of sections
// 4. Change the IMAGE_SIZE
// 5. Add bytes to memory code-cave
