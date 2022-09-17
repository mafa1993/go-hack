

1. 使用debug/pe标准库进行解析
2. 使用Reader对象对PE文件内容进行解析

# PE文件结构

1. DOSheader 包含签名（0x5a4d）peheader(0x3c指向0x50 0x45 0x00 0x00)
2. dos stub
3. coff file header
4. optional header
5. section table

```
// pe头
type FileHeader struct {
	Machine              uint16
	NumberOfSections     uint16  // 分区数
	TimeDateStamp        uint32
	PointerToSymbolTable uint32
	NumberOfSymbols      uint32
	SizeOfOptionalHeader uint16
	Characteristics      uint16
}
```
1. 如果需要增加新分区，插入后门代码，需要修改这里的分区数属性
