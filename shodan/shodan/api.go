package shodan

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// api 获取余量的

type APIInfo struct {
	QueryCredits int    `json:"query_credits`
	ScanCredits  int    `json:"scan_credits"`
	Telnet       bool   `json:"telnet`
	Plan         string `json:"plan"`
	HTTPS        bool   `json:"https"`
	Unlocked     bool   `json:"unlocked"`
}

func (client *Client) APIInfo() (*APIInfo, error) {
	apikey := client.apiKey
	url := fmt.Sprintf("%s/api-info?key=%s", BASEURL, apikey)
	fmt.Println("url", url)

	resp, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return &APIInfo{}, err
	}

	// s, _ := ioutil.ReadAll(resp.Body)
	// fmt.Printf("%s", s)
	var apiInfo APIInfo
	err = json.NewDecoder(resp.Body).Decode(&apiInfo)
	//fmt.Println(apiInfo)
	//fmt.Printf("%+v", apiInfo)
	if err != nil {
		return &APIInfo{}, err
	}
	defer resp.Body.Close()
	return &apiInfo, nil
}
