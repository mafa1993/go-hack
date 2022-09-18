package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
)

// 头
type Header struct {
	Header uint64
}

// 数据块
type Chunk struct {
	Size uint32
	Type uint32
	Data []byte
	CRC  uint32
}

type CmdLineOpts struct {
	Input    string
	Output   string
	Meta     bool
	Suppress bool
	Offset   string
	Inject   bool
	Payload  string
	Type     string
	Encode   bool
	Decode   bool
	Key      string
}

//MetaChunk inherits a Chunk struct
type MetaChunk struct {
	Chk    Chunk
	Offset int64
}

var (
	flags = pflag.FlagSet{SortFlags: false}
	opts  CmdLineOpts
	png   MetaChunk
)

func init() {
	flags.StringVarP(&opts.Input, "input", "i", "", "Path to the original image file")
	flags.StringVarP(&opts.Output, "output", "o", "", "Path to output the new image file")
	flags.BoolVarP(&opts.Meta, "meta", "m", false, "Display the actual image meta details")
	flags.BoolVarP(&opts.Suppress, "suppress", "s", false, "Suppress the chunk hex data (can be large)")
	flags.StringVar(&opts.Offset, "offset", "", "The offset location to initiate data injection")
	flags.BoolVar(&opts.Inject, "inject", false, "Enable this to inject data at the offset location specified")
	flags.StringVar(&opts.Payload, "payload", "", "Payload is data that will be read as a byte stream")
	flags.StringVar(&opts.Type, "type", "rNDm", "Type is the name of the Chunk header to inject")
	flags.StringVar(&opts.Key, "key", "", "The enryption key for payload")
	flags.BoolVar(&opts.Encode, "encode", false, "XOR encode the payload")
	flags.BoolVar(&opts.Decode, "decode", false, "XOR decode the payload")
	flags.Lookup("type").NoOptDefVal = "rNDm"
	flags.Usage = usage
	flags.Parse(os.Args[1:])

	if flags.NFlag() == 0 {
		flags.PrintDefaults()
		os.Exit(1)
	}
	if opts.Input == "" {
		log.Fatal("Fatal: The --input flag is required")
	}
	if opts.Offset != "" {
		byteOffset, _ := strconv.ParseInt(opts.Offset, 0, 64)
		opts.Offset = strconv.FormatInt(byteOffset, 10)
	}
	if opts.Suppress && (opts.Meta == false) {
		log.Fatal("Fatal: The --meta flag is required when using --suppress")
	}
	if opts.Meta && (opts.Offset != "") {
		log.Fatal("Fatal: The --meta flag is mutually exclusive with --offset")
	}
	if opts.Inject && (opts.Offset == "") {
		log.Fatal("Fatal: The --offset flag is required when using --inject")
	}
	if opts.Inject && (opts.Payload == "") {
		log.Fatal("Fatal: The --payload flag is required when using --inject")
	}
	if opts.Inject && opts.Key == "" {
		fmt.Println("Warning: No key provided. Payload will not be encrypted")
	}
	if opts.Encode && opts.Key == "" {
		log.Fatal("Fatal: The --encode flag requires a --key value")
	}
}

func usage() {
	fmt.Fprintf(os.Stderr, "Example Usage: %s -i in.png -o out.png --inject --offset 0x85258 --payload 1234\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Example Encode Usage: %s -i in.png -o encode.png --inject --offset 0x85258 --payload 1234 --encode --key secret\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Example Decode Usage: %s -i encode.png -o decode.png --offset 0x85258 --decode --key secret\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Flags: %s {OPTION]...\n", os.Args[0])
	flags.PrintDefaults()
	os.Exit(0)
}

func main() {
	dat, err := os.Open(opts.Input)
	defer dat.Close()
	bReader, err := PreProcessImage(dat)
	if err != nil {
		log.Fatal(err)
	}
	png.ProcessImage(bReader, &opts)
}

// 从文件句柄读到缓冲区
func PreProcessImage(dat *os.File) (*bytes.Reader, error) {
	state, err := dat.Stat() // 获取文件的state信息

	if err != nil {
		return nil, err
	}

	var size = state.Size() // 获取文件大小
	b := make([]byte, size) // 创建一个缓冲区用来读取文件

	buf := bufio.NewReader(dat) // 对dat创建一个Reader对象

	_, err = buf.Read(b) // 将dat copy到b中  返回大小和错误

	if err != nil {
		return nil, err
	}

	breader := bytes.NewReader(b) // 对b创建一个Reader对象

	return breader, nil
}

func (mc *MetaChunk) validate(b *bytes.Reader) {
	var header Header

	// 将前8个字节写到header.Header中
	if err := binary.Read(b, binary.BigEndian, &header.Header); err != nil {
		log.Fatalln(err)
	}

	bArr := make([]byte, 8)
	binary.BigEndian.PutUint64(bArr, header.Header) // header.Header转换为[]byte

	// 验证2-4字节值是否为PNG
	if string(bArr[1:4]) != "PNG" {
		log.Fatal("文件类型错误")
	}
	fmt.Println("文件类型验证通过")
}

// 块解析，size、type、data、crc循环
func (mc *MetaChunk) ProcessingImage(b *bytes.Reader, c *CmdLineOpts) {

	if (c.Offset != "") && (c.Encode == false && c.Decode == false) {
		var m MetaChunk
		m.Chk.Data = []byte(c.Payload)
		m.Chk.Type = m.strToInt(c.Type)
		m.Chk.Size = m.createChunkSize()
		m.Chk.CRC = m.createChunkCRC()
		bm := m.marshalData()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", []byte(c.Payload))
		fmt.Printf("Payload: % X\n", m.Chk.Data)
		WriteData(b, c, bmb)
	}
	if (c.Offset != "") && c.Encode {
		var m MetaChunk
		m.Chk.Data = XorEncode([]byte(c.Payload), c.Key)
		m.Chk.Type = m.strToInt(c.Type)
		m.Chk.Size = m.createChunkSize()
		m.Chk.CRC = m.createChunkCRC()
		bm := m.marshalData()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", []byte(c.Payload))
		fmt.Printf("Payload Encode: % X\n", m.Chk.Data)
		WriteData(b, c, bmb)
	}
	if (c.Offset != "") && c.Decode {
		var m MetaChunk
		offset, _ := strconv.ParseInt(c.Offset, 10, 64)
		b.Seek(offset, 0)
		m.readChunk(b)
		origData := m.Chk.Data
		m.Chk.Data = XorDecode(m.Chk.Data, c.Key)
		m.Chk.CRC = m.createChunkCRC()
		bm := m.marshalData()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", origData)
		fmt.Printf("Payload Decode: % X\n", m.Chk.Data)
		WriteData(b, c, bmb)
	}
	if c.Meta {
		count := 1 //Start at 1 because 0 is reserved for magic byte
		var chunkType string
		for chunkType != "IEND" {
			mc.getOffset(b)
			mc.readChunk(b)
			fmt.Println("第", count, "块")
			fmt.Printf("Chunk Offset: %#02x\n", mc.Offset)
			fmt.Printf("Chunk Length: %s bytes\n", strconv.Itoa(int(mc.Chk.Size)))
			fmt.Printf("Chunk Type: %s\n", mc.chunkTypeToString())
			fmt.Printf("Chunk Importance: %s\n", mc.checkCritType())
			if c.Suppress == false {
				fmt.Printf("Chunk Data: %#x\n", mc.Chk.Data)
			} else if c.Suppress {
				fmt.Printf("Chunk Data: %s\n", "Suppressed")
			}
			fmt.Printf("Chunk CRC: %x\n", mc.Chk.CRC)
			chunkType = mc.chunkTypeToString()
			count++
		}
	}

}

func (mc *MetaChunk) getOffset(b *bytes.Reader) {
	offset, _ := b.Seek(0, 1) // 返回偏移量，1代表从当前位置
	mc.Offset = offset
}

// 块读取方法
func (mc *MetaChunk) readChunk(b *bytes.Reader) {
	mc.readChunkSize(b)
	mc.readChunkType(b)
	mc.readChunkBytes(b, mc.Chk.Size) // 数据读取
	mc.readChunkCRC(b)
}

// 获取块大小
func (mc *MetaChunk) readChunkSize(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Size); err != nil {
		log.Fatalln(err)
	}
}

// 获取块类型
func (mc *MetaChunk) readChunkType(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Type); err != nil {
		log.Fatalln(err)
	}
}

func (mc *MetaChunk) readChunkBytes(b *bytes.Reader, cLen uint32) {
	mc.Chk.Data = make([]byte, cLen)
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.Data); err != nil {
		log.Fatalln(err)
	}
}

func (mc *MetaChunk) readChunkCRC(b *bytes.Reader) {
	if err := binary.Read(b, binary.BigEndian, &mc.Chk.CRC); err != nil {
		log.Fatalln(err)
	}
}

func (mc *MetaChunk) strToInt(s string) uint32 {
	t := []byte(s)
	return binary.BigEndian.Uint32(t)
}

func (mc *MetaChunk) createChunk() uint32 {
	return uint32(len(mc.Chk.Data)) // 转成uint32返回数据长度
}

// 校验type和data值
func (mc *MetaChunk) createChunkCRC() uint32 {
	// 创建一个buffer，将type和data放入
	bytesMSB := new(bytes.Buffer)

	err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Type)
	if err != nil {
		log.Fatalln(err)
	}

	err = binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Data)
	if err != nil {
		log.Fatalln(err)
	}

	return crc32.ChecksumIEEE(bytesMSB.Bytes()) // 计算校验和
}

func (mc *MetaChunk) marshalData() *bytes.Buffer {
	bytesMSB := new(bytes.Buffer)

	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Size); err != nil {
		log.Fatalln(err)
	}

	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Type); err != nil {
		log.Fatalln(err)
	}

	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.Data); err != nil {
		log.Fatalln(err)
	}

	if err := binary.Write(bytesMSB, binary.BigEndian, mc.Chk.CRC); err != nil {
		log.Fatalln(err)
	}

	return bytesMSB
}

// @param r 原始数据
// @param b 新组成的数据
func WriteData(r *bytes.Reader, c *CmdLineOpts, b []byte) {
	offset, _ := strconv.ParseInt(c.Offset, 10, 64) // 字符串转int，10进制的字符串 返回64位, offset为头的偏移量 一般64
	w, err := os.Create(c.Output)                   // 创建输出的文件
	if err != nil {
		log.Fatalln(err)
	}

	defer w.Close()

	r.Seek(0, 0)                   // 重置指针
	var buf = make([]byte, offset) // 创建存储头的
	r.Read(buf)                    // 读取头
	w.Write(buf)                   // 头写入到w
	w.Write(b)                     // payload写入
	_, err = io.Copy(w, r)         // 剩余的数据写到新文件
	if err != nil {
		log.Fatalln(err)
	}

}

//ProcessImage is the wrapper to parse PNG bytes
func (mc *MetaChunk) ProcessImage(b *bytes.Reader, c *CmdLineOpts) {
	mc.validate(b)
	if (c.Offset != "") && (c.Encode == false && c.Decode == false) {
		var m MetaChunk
		m.Chk.Data = []byte(c.Payload)
		m.Chk.Type = m.strToInt(c.Type)
		m.Chk.Size = m.createChunkSize()
		m.Chk.CRC = m.createChunkCRC()
		bm := m.marshalData()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", []byte(c.Payload))
		fmt.Printf("Payload: % X\n", m.Chk.Data)
		WriteData(b, c, bmb)
	}
	if (c.Offset != "") && c.Encode {
		var m MetaChunk
		m.Chk.Data = XorEncode([]byte(c.Payload), c.Key)
		m.Chk.Type = m.strToInt(c.Type)
		m.Chk.Size = m.createChunkSize()
		m.Chk.CRC = m.createChunkCRC()
		bm := m.marshalData()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", []byte(c.Payload))
		fmt.Printf("Payload Encode: % X\n", m.Chk.Data)
		WriteData(b, c, bmb)
	}
	if (c.Offset != "") && c.Decode {
		var m MetaChunk
		offset, _ := strconv.ParseInt(c.Offset, 10, 64)
		b.Seek(offset, 0)
		m.readChunk(b)
		origData := m.Chk.Data
		m.Chk.Data = XorDecode(m.Chk.Data, c.Key)
		m.Chk.CRC = m.createChunkCRC()
		bm := m.marshalData()
		bmb := bm.Bytes()
		fmt.Printf("Payload Original: % X\n", origData)
		fmt.Printf("Payload Decode: % X\n", m.Chk.Data)
		WriteData(b, c, bmb)
	}
	if c.Meta {
		count := 1 //Start at 1 because 0 is reserved for magic byte
		var chunkType string
		for chunkType != "IEND" {
			mc.getOffset(b)
			mc.readChunk(b)
			fmt.Println("---- Chunk # " + strconv.Itoa(count) + " ----")
			fmt.Printf("Chunk Offset: %#02x\n", mc.Offset)
			fmt.Printf("Chunk Length: %s bytes\n", strconv.Itoa(int(mc.Chk.Size)))
			fmt.Printf("Chunk Type: %s\n", mc.chunkTypeToString())
			fmt.Printf("Chunk Importance: %s\n", mc.checkCritType())
			if c.Suppress == false {
				fmt.Printf("Chunk Data: %#x\n", mc.Chk.Data)
			} else if c.Suppress {
				fmt.Printf("Chunk Data: %s\n", "Suppressed")
			}
			fmt.Printf("Chunk CRC: %x\n", mc.Chk.CRC)
			chunkType = mc.chunkTypeToString()
			count++
			if count == 3 {
				break
			}
		}
	}
}

func (mc *MetaChunk) createChunkSize() uint32 {
	return uint32(len(mc.Chk.Data))
}

func (mc *MetaChunk) chunkTypeToString() string {
	h := fmt.Sprintf("%x", mc.Chk.Type)
	decoded, _ := hex.DecodeString(h)
	result := fmt.Sprintf("%s", decoded)
	return result
}

func (mc *MetaChunk) checkCritType() string {
	fChar := string([]rune(mc.chunkTypeToString())[0])
	if fChar == strings.ToUpper(fChar) {
		return "Critical"
	}
	return "Ancillary"
}
