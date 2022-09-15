package main

import (
	"fmt"
	"goplugin/scanner"
	"log"
	"net/http"
)

// 构建插件的主入口

// 构造字典
var User = []string{
	"admin", "tomacat",
}
var Password = []string{
	"admin", "123456",
}

type TomacatChecker struct{}

// 实现scanner接口
func (c TomacatChecker) Check(host string, port uint64) *scanner.Result {
	var (
		resp   *http.Response
		err    error
		url    string
		res    *scanner.Result
		client *http.Client
		req    *http.Request
	)

	log.Println("检查tomcat401")

	res = new(scanner.Result) // new 创建一个空结构体  返回指针

	url = fmt.Sprintf("http://%s:%d/manager", host, port) // 构建请求的url

	// 发送head请求，进行测试
	if resp, err = http.Head(url); err != nil {
		log.Fatalln("发送head请求出错", err)
		return res
	}

	// 检查身份认证类型
	// 如果状态码部位401获取响应头没有www-Authenticate，报错返回
	// resp.Header.Get 获取响应头
	if resp.StatusCode != http.StatusUnauthorized || resp.Header.Get("www-Authenticate") == "" {
		log.Fatalln("不需要认证")
		return res
	}

	fmt.Println("进行密码猜解")
	client = new(http.Client) // client = &http.Client{}
	if req, err = http.NewRequest("GET", url, nil); err != nil {
		log.Println("建立请求失败")
		return res
	}

	for _, user := range User {
		for _, pass := range Password {
			req.SetBasicAuth(user, pass) // 设置401 请求头

			if resp, err = client.Do(req); err != nil {
				log.Println("请求出错", err)
				continue
			}

			if resp.StatusCode == http.StatusOK {
				res.Vulnerable = true
				res.Details = fmt.Sprintf("user:%s,pass:%s", user, pass)
				return res
			}

		}

	}

	// 到最后没有找到正确的账密
	return res

}

func New() scanner.Checker {
	return new(TomacatChecker)
}

//go build -buildmode=plugin -o ./plugins/
