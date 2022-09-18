
对png文件进行解析，并进行修改，实现将shellcode放入png文件中，在使用时，从png图片中获取

# 准备工作

1. http://www.libpng.org/pub/png/spec/1.2/PNG-Structure.html  png结构


# png文件结构

1. 前8个字节位header头，结尾均为CRLF
2. SIZE 4个字节 定义了随后的data长度
3. TYPE 4个字节  IHDR | IDAT 定义类型，枚举
4. DATA 任意数量
5. CRC 4个字节  对type和data的数据进行校验，crc-32校验的


//go run .\main.go .\utils.go -i "./php.png" -o "a.png"  --inject --offset 0x85258 --payload 1233333