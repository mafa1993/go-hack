package metadata

import (
	"archive/zip"
	"encoding/xml"
	"log"
	"strings"
)

type OfficeCoreProperty struct {
	XMLName        xml.Name `xml:"coreProperties`
	Creator        string   `xml:"creator"`        // 创建者
	LastModifiedBy string   `xml:"lastModifiedBy"` // 最后修改人
}

type OfficeAppProperty struct {
	XMLName     xml.Name `xml:"Properties"`
	Application string   `xml:"Application"`
	Company     string   `xml:"Company"`
	Version     string   `xml:"AppVersion"`
}

var OfficeVersions = map[string]string{
	"16": "2016",
	"15": "2013",
	"14": "2010",
	"12": "2007",
	"11": "2003",
}

// 将版本转换为可读的版本 , 如果没有匹配上返回unknown
func (p OfficeAppProperty) GetMajorVersion() string {
	version := strings.Split(p.Version, ".") // 截取到点也可以
	rlt, ok := OfficeVersions[version[0]]
	if !ok {
		return "unknow"
	}
	return rlt
}

/**
 * 根据结构体解析文件
 * @params r 从网络中获取到的doc等文件，对其转换成zip.Reader进行读取
 */
func NewProperties(r *zip.Reader) (*OfficeCoreProperty, *OfficeAppProperty, error) {
	var (
		coreProps OfficeCoreProperty
		appProps  OfficeAppProperty
	)
	for _, v := range r.File {
		switch v.Name {
		// core.xml 解析
		case "docProps/core.xml":
			if err := Process(v, &coreProps); err != nil {
				return nil, nil, err
			}
		// app.xml 机械
		case "docProps/app.xml":
			if err := Process(v, &appProps); err != nil {
				return nil, nil, err
			}
		}
	}
	return &coreProps, &appProps, nil

}

// 文件解析，
// f 文件名
// rlt 解析的结果,这里interface可以接收任意类型，这里应该为指针类型
func Process(file *zip.File, rlt interface{}) error {
	f, err := file.Open()
	defer f.Close() // f 为io.ReaderCloser
	if err != nil {
		log.Println(err)
		return err
	}
	err = xml.NewDecoder(f).Decode(rlt)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
