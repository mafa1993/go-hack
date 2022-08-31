package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 查询具体内容的函数

// https://api.shodan.io/shodan/host/search?key=?&query=?&facets=?

// location 数据
type HostLocation struct {
	City         string  `json:"city"`
	RegionCode   string  `json:"region_code"`
	AreaCode     int     `json:"area_code"`
	Longitude    float32 `json:"longitude"`
	CountryCode3 string  `json:"country_code3"`
	CountryName  string  `json:"country_name"`
	PostalCode   int     `json:"postal_code"`
	DMACode      string  `json:"dma_code"`
	Latitude     float32 `json:"latitude"`
}

// matches 每一个元素
type Host struct {
	OS        string       `json:"os"`
	Location  HostLocation `json:"location"`
	Timestamp string       `json:"timestamp"`
	ISP       string       `json:"isp"`
	ASN       string       `json:"asn"`
	Hostname  []string     `json:"hostname"`
	IP        int64        `json:"ip"`
	Domain    []string     `json:"domain"`
	Org       string       `json:"org"`
	Data      string       `json:"data"`
	Port      int          `json:"port"`
	IPString  string       `json:"ip_str"`
}

// 解析matches数组
type HostSearch struct {
	Matches []Host `json:"matches"`
}

// @params q 查询的内容
// 此函数用于查询
func (client *Client) HostSearch(q string) (*HostSearch, error) {
	url := fmt.Sprintf("%s/shodan/host/search?key=%s&query=%s", BASEURL, client.apiKey, q)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Println(err)
		return &HostSearch{}, err
	}
	var c http.Client
	resp, err := c.Do(req)
	if err != nil {
		fmt.Println(err)
		return &HostSearch{}, err
	}

	defer resp.Body.Close()

	var rtn HostSearch
	if err := json.NewDecoder(resp.Body).Decode(&rtn); err != nil {
		fmt.Println(err)
		return &HostSearch{}, err
	}
	fmt.Printf("%+v", rtn)

	return &rtn, nil
}
