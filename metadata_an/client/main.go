package main

import (
	"archive/zip"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"metadata_an/metadata"
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func init() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

const SEARCH_URL = "http://cn.bing.com/search?q=%s"

// 从命令行传入参数，进行bing搜索
func main() {
	//site := flag.String("site", "nytimes.com", "网址")
	filetype := flag.String("filetype", "docx", "文件类型")
	flag.Parse()
	// fmt.Println(flag.NFlag())
	// if flag.NFlag() < 2 {
	// 	fmt.Println("参数错误")
	// 	panic(1)
	// }

	//q := fmt.Sprintf("site:%s && filetype:%s && instreamset:(url title):%s", *site, *filetype, *filetype)
	q := fmt.Sprintf("filetype:%s", *filetype)
	search := fmt.Sprintf(SEARCH_URL, q)
	//search = SEARCH_URL

	//search = "https://www.baidu.com/s?wd=filetype%3Adocx"
	doc, err := goquery.NewDocument(search) // 会发送Get请求

	if err != nil {
		log.Println(err)
		panic(err)
	}

	s := "li"                 // 选择器
	doc.Find(s).Each(handler) // 根据s选择器获取每个元素，调用handler函数
}

func handler(i int, s *goquery.Selection) {

	url, ok := s.Find("a").Attr("href")
	fmt.Println(i, url)
	if !ok {
		return
	}

	res, err := http.Get(url) // 发送请求获取文档

	if err != nil {
		log.Println(err)
		return
	}
	defer res.Body.Close()

	// 读取全部数据
	data, _ := ioutil.ReadAll(res.Body)

	// 创建一个zip的reader 读取获取的信息，利用zip格式读取下载的内容，建立的缓冲区大小为全部内容的缓冲区
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))

	if err != nil {
		log.Println(err)
		return
	}

	cp, ap, err := metadata.NewProperties(r)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("创建者：%25s,最后修改：%25s-文件类型：%s 文件版本：%s\n", cp.Creator, cp.LastModifiedBy, ap.Application, ap.GetMajorVersion())
	return
}
